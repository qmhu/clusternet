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

package socket

import (
	"context"
	"fmt"
	"github.com/clusternet/clusternet/pkg/apis/federations/v1alpha1"
	"k8s.io/client-go/rest"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	registryrest "k8s.io/apiserver/pkg/registry/rest"

	"github.com/clusternet/clusternet/pkg/federation"
)

const (
	category = "clusternet"
)

// REST implements a RESTStorage for federation API
type REST struct {
	manager        *federation.Manager
}

func (r *REST) ShortNames() []string {
	return []string{"govern"}
}

func (r *REST) NamespaceScoped() bool {
	return false
}

func (r *REST) Categories() []string {
	return []string{category}
}

func (r *REST) New() runtime.Object {
	return &v1alpha1.Govern{}
}

// TODO: constraint govern methods
var governMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// ConnectMethods returns the list of HTTP methods that can be federated
func (r *REST) ConnectMethods() []string {
	return governMethods
}

// NewConnectOptions returns versioned resource that represents federated parameters
func (r *REST) NewConnectOptions() (runtime.Object, bool, string) {
	return &v1alpha1.Govern{}, true, ""
}

// Connect returns a handler for the websocket connection
func (r *REST) Connect(ctx context.Context, id string, opts runtime.Object, responder registryrest.Responder) (http.Handler, error) {
	govern, ok := opts.(*v1alpha1.Govern)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", opts)
	}
	return r.manager.HandleConnection(ctx, id, govern, responder)
}

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(config *rest.Config) *REST {
	return &REST{
		manager: federation.NewManager(config),
	}
}

var _ registryrest.CategoriesProvider = &REST{}
var _ registryrest.ShortNamesProvider = &REST{}
var _ registryrest.Connecter = &REST{}
