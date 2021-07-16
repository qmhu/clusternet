package handler

import (
	"context"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/proxy"
	registryrest "k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

type ProxyHandler struct {
	config   *rest.Config
	proxyUrl string
}

func NewProxyHandler(config *rest.Config, proxyUrl string) *ProxyHandler {
	return &ProxyHandler{
		config:   config,
		proxyUrl: proxyUrl,
	}
}

func (h *ProxyHandler) HandleConnection(ctx context.Context, w http.ResponseWriter, req *http.Request, responder registryrest.Responder) error {
	reverseProxy, hubApiserverUrl, err := NewReverseProxyForHubApiserver(h.config)
	if err != nil {
		return apierrors.NewInternalError(err)
	}

	reverseProxy.Director = func(req *http.Request) {
		req.URL.Path = h.proxyUrl
		req.URL.Scheme = hubApiserverUrl.Scheme
		req.URL.Host = hubApiserverUrl.Host

		klog.V(4).Infof("Forward federation request, requesturi %s", req.URL.String())
	}

	reverseProxy.ErrorHandler = proxy.NewErrorResponder(responder).Error
	reverseProxy.ServeHTTP(w, req)
	return nil
}
