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

package constants

import "time"

const (
	ClusternetSystemNamespace = "clusternet-system"
	ManifestInstanceIDLength  = 5
	ManifestRequestInterval   = 1 * time.Second
	ManifestRequestTimeout    = 5 * time.Second

	FeedUidLabel        string = "apps.clusternet.io/feed.uid"
	FeedApiVersionLabel string = "apps.clusternet.io/feed.apiversion"
	FeedKindLabel       string = "apps.clusternet.io/feed.kind"
	FeedNamespaceLabel  string = "apps.clusternet.io/feed.namespace"
	FeedNameLabel       string = "apps.clusternet.io/feed.name"

	RequestPathDeclarations   string = "declarations"
)
