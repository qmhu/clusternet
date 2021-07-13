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

package apiserver

import (
	"github.com/clusternet/clusternet/pkg/apis/federations"
	"github.com/clusternet/clusternet/pkg/apis/proxies"
	"github.com/clusternet/clusternet/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	registryrest "k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"

	federationsInstall "github.com/clusternet/clusternet/pkg/apis/federations/install"
	proxiesInstall "github.com/clusternet/clusternet/pkg/apis/proxies/install"

	"github.com/clusternet/clusternet/pkg/exchanger"
	clusterInformers "github.com/clusternet/clusternet/pkg/generated/informers/externalversions/clusters/v1beta1"
	governstorage "github.com/clusternet/clusternet/pkg/registry/federation/govern"
	socketstorage "github.com/clusternet/clusternet/pkg/registry/proxies/socket"
	"github.com/clusternet/clusternet/pkg/registry/proxies/socket/subresources"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme)
	// ParameterCodec handles versioning of objects that are converted to query parameters.
	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

func init() {
	proxiesInstall.Install(Scheme)
	federationsInstall.Install(Scheme)

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	// Place you custom config here.
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

// HubAPIServer contains state for a master/api server.
type HubAPIServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}

	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return CompletedConfig{&c}
}

// New returns a new instance of HubAPIServer from the given config.
func (c completedConfig) New(tunnelLogging, socketConnection bool, extraHeaderPrefixes []string,
	mclsInformer clusterInformers.ManagedClusterInformer, CoreAPIKubeconfigPath string) (*HubAPIServer, error) {
	genericServer, err := c.GenericConfig.New("clusternet-hub", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	config, err := utils.LoadsKubeConfig(CoreAPIKubeconfigPath, 10)
	if err != nil {
		return nil, err
	}

	s := &HubAPIServer{
		GenericAPIServer: genericServer,
	}

	apiGroupInfoProxies := genericapiserver.NewDefaultAPIGroupInfo(proxies.GroupName, Scheme, ParameterCodec, Codecs)

	var ec *exchanger.Exchanger
	if socketConnection {
		ec = exchanger.NewExchanger(tunnelLogging, mclsInformer)
	}
	v1alpha1storageProxies := map[string]registryrest.Storage{}
	v1alpha1storageProxies["sockets"] = socketstorage.NewREST(socketConnection, ec)
	v1alpha1storageProxies["sockets/proxy"] = subresources.NewProxyREST(socketConnection, ec, extraHeaderPrefixes)

	apiGroupInfoProxies.VersionedResourcesStorageMap["v1alpha1"] = v1alpha1storageProxies

	apiGroupInfoFeds := genericapiserver.NewDefaultAPIGroupInfo(federations.GroupName, Scheme, ParameterCodec, Codecs)

	v1alpha1storageFeds := map[string]registryrest.Storage{}
	v1alpha1storageFeds["governs"] = governstorage.NewREST(config)
	apiGroupInfoFeds.VersionedResourcesStorageMap["v1alpha1"] = v1alpha1storageFeds


	if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfoProxies); err != nil {
		return nil, err
	}

	if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfoFeds); err != nil {
		return nil, err
	}

	return s, nil
}
