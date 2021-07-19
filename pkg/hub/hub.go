/*
Copyright 2021 The Clusternet Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hub

import (
	"context"
	"github.com/clusternet/clusternet/pkg/federation"
	"time"

	genericapiserver "k8s.io/apiserver/pkg/server"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	kubeInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	"github.com/clusternet/clusternet/pkg/features"
	clusternetClientSet "github.com/clusternet/clusternet/pkg/generated/clientset/versioned"
	informers "github.com/clusternet/clusternet/pkg/generated/informers/externalversions"
	"github.com/clusternet/clusternet/pkg/hub/approver"
	"github.com/clusternet/clusternet/pkg/hub/deployer"
	"github.com/clusternet/clusternet/pkg/hub/options"
	"github.com/clusternet/clusternet/pkg/utils"
)

const (
	// default resync time
	DefaultResync = time.Hour * 12
	// default number of threads
	DefaultThreadiness = 2
)

// Hub defines configuration for clusternet-hub
type Hub struct {
	ctx                       context.Context
	options                   *options.HubServerOptions
	crrApprover               *approver.CRRApprover
	clusternetInformerFactory informers.SharedInformerFactory
	kubeInformerFactory       kubeInformers.SharedInformerFactory
	kubeclient                *kubernetes.Clientset
	clusternetclient          *clusternetClientSet.Clientset
	deployer                  *deployer.Deployer
	fedManager                *federation.Server
	socketConnection          bool
	deployerEnabled           bool
	resourceAsAppsEnabled	  bool
}

// NewHub returns a new Hub.
func NewHub(ctx context.Context, opts *options.HubServerOptions) (*Hub, error) {
	socketConnection := utilfeature.DefaultFeatureGate.Enabled(features.SocketConnection)
	deployerEnabled := utilfeature.DefaultFeatureGate.Enabled(features.Deployer)
	resourceAsAppsEnabled := utilfeature.DefaultFeatureGate.Enabled(features.ResourceAsApps)

	config, err := utils.LoadsKubeConfig(opts.RecommendedOptions.CoreAPI.CoreAPIKubeconfigPath, 10)
	if err != nil {
		return nil, err
	}

	// creating the clientset
	kubeclient := kubernetes.NewForConfigOrDie(config)
	clusternetclient := clusternetClientSet.NewForConfigOrDie(config)

	// creates the informer factory
	kubeInformerFactory := kubeInformers.NewSharedInformerFactory(kubeclient, DefaultResync)
	clusternetInformerFactory := informers.NewSharedInformerFactory(clusternetclient, DefaultResync)
	approver, err := approver.NewCRRApprover(ctx, kubeclient, clusternetclient, clusternetInformerFactory,
		kubeInformerFactory, socketConnection)
	if err != nil {
		return nil, err
	}

	deployer, err := deployer.NewDeployer(ctx, kubeclient, clusternetclient, clusternetInformerFactory, kubeInformerFactory)
	if err != nil {
		return nil, err
	}

	fedManager := federation.NewServer(config, kubeclient, clusternetclient)

	// add informers
	kubeInformerFactory.Core().V1().Namespaces().Informer()
	kubeInformerFactory.Core().V1().ServiceAccounts().Informer()
	kubeInformerFactory.Core().V1().Secrets().Informer()
	clusternetInformerFactory.Clusters().V1beta1().ClusterRegistrationRequests().Informer()
	clusternetInformerFactory.Clusters().V1beta1().ManagedClusters().Informer()
	if deployerEnabled {
		clusternetInformerFactory.Apps().V1alpha1().Subscriptions().Informer()
		clusternetInformerFactory.Apps().V1alpha1().HelmCharts().Informer()
		clusternetInformerFactory.Apps().V1alpha1().Descriptions().Informer()
		clusternetInformerFactory.Apps().V1alpha1().HelmReleases().Informer()
	}

	hub := &Hub{
		ctx:                       ctx,
		crrApprover:               approver,
		options:                   opts,
		kubeclient:                kubeclient,
		clusternetclient:          clusternetclient,
		clusternetInformerFactory: clusternetInformerFactory,
		kubeInformerFactory:       kubeInformerFactory,
		socketConnection:          socketConnection,
		deployer:                  deployer,
		deployerEnabled:           deployerEnabled,
		fedManager:                fedManager,
		resourceAsAppsEnabled:     resourceAsAppsEnabled,
	}

	// Start the informer factories to begin populating the informer caches
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(ctx.Done())
	clusternetInformerFactory.Start(ctx.Done())

	return hub, nil
}

func (hub *Hub) Run() error {
	go func() {
		hub.crrApprover.Run(DefaultThreadiness)
	}()

	if hub.deployerEnabled {
		go func() {
			hub.deployer.Run(DefaultThreadiness)
		}()
	}

	err := hub.RunAPIServer()
	if err != nil {
		return err
	}

	return nil
}

// RunAPIServer starts a new HubAPIServer given HubServerOptions
func (hub *Hub) RunAPIServer() error {
	klog.Info("starting Clusternet Hub APIServer ...")
	config, err := hub.options.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New(hub.options.TunnelLogging, hub.socketConnection,
		hub.options.RecommendedOptions.Authentication.RequestHeader.ExtraHeaderPrefixes,
		hub.clusternetInformerFactory.Clusters().V1beta1().ManagedClusters(),
		hub.fedManager)
	if err != nil {
		return err
	}

	hub.fedManager.ApiserverConfig = config.GenericConfig

	server.GenericAPIServer.AddPostStartHookOrDie("start-clusternet-hub-apiserver-informers", func(context genericapiserver.PostStartHookContext) error {
		config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
		// no need to start LoopbackSharedInformerFactory since we don't store anything in this apiserver
		// hub.options.LoopbackSharedInformerFactory.Start(context.StopCh)
		return nil
	})

	return server.GenericAPIServer.PrepareRun().Run(hub.ctx.Done())
}
