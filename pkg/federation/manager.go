package federation

import (
	"context"
	"fmt"
	types "github.com/clusternet/clusternet/pkg/apis/federations/v1alpha1"
	registryrest "k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/rest"
	"net/http"
)

type Manager struct {
	config *rest.Config
}

var (
	urlPrefix = fmt.Sprintf("/apis/%s/governs/", types.SchemeGroupVersion.String())
)

func NewManager(config *rest.Config) *Manager {
	return &Manager{
		config: config,
	}
}

func (m *Manager) HandleConnection(ctx context.Context, id string, opts *types.Govern, responder registryrest.Responder) (http.Handler, error)  {
	proxy := NewReverseProxyHandler(m.config, ctx, id, opts, responder)
	return proxy, nil
}

