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

package features

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/component-base/featuregate"
)

const (
	// alpha: v0.1.0
	//
	// Setup/Serve a WebSocket connection.
	SocketConnection featuregate.Feature = "SocketConnection"

	// alpha: v0.2.0
	//
	// Allow to deploy applications from parent cluster.
	// Mainly for security concerns of every child cluster.
	// If a child cluster has disabled AppPusher, the parent cluster won't deploy applications with Push or Dual mode.
	AppPusher featuregate.Feature = "AppPusher"

	// alpha: v0.2.0
	//
	// Works as a deployer that help distribute kinds of resources to a group of clusters
	Deployer featuregate.Feature = "Deployer"

	// alpha: v0.3.0
	//
	// Accept k8s resource and transform into manifest.
	// Manifests can be distributed to a group of clusters as a App.
	ResourceAsApps featuregate.Feature = "ResourceAsApps"

	// TODO
)

func init() {
	runtime.Must(utilfeature.DefaultMutableFeatureGate.Add(defaultClusternetFeatureGates))
}

// defaultClusternetFeatureGates consists of all known Kubernetes-specific and clusternet feature keys.
// To add a new feature, define a key for it above and add it here. The features will be
// available throughout Clusternet binaries.
var defaultClusternetFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	SocketConnection: {Default: false, PreRelease: featuregate.Alpha, LockToDefault: false},
	AppPusher:        {Default: false, PreRelease: featuregate.Alpha, LockToDefault: false},
	Deployer:         {Default: false, PreRelease: featuregate.Alpha, LockToDefault: false},
	ResourceAsApps:   {Default: false, PreRelease: featuregate.Alpha, LockToDefault: false},
}
