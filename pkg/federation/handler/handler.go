package handler

import (
	"context"
	"k8s.io/client-go/rest"
	"net/http"
	"net/http/httputil"
	"net/url"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	registryrest "k8s.io/apiserver/pkg/registry/rest"
)

type RequestHandler interface {
	HandleConnection(ctx context.Context, w http.ResponseWriter, req *http.Request, responder registryrest.Responder) error
}

type DummyHandler struct{}

func NewDummyHandler() *DummyHandler {
	return &DummyHandler{}
}

func (h *DummyHandler) HandleConnection(ctx context.Context, w http.ResponseWriter, req *http.Request, responder registryrest.Responder) error {
	return apierrors.NewBadRequest("not support this kind of request yet.")
}

func DryRun(req *http.Request) {
	q := req.URL.Query()
	q.Add("dryRun", "All")
	req.URL.RawQuery = q.Encode()
}

func NewReverseProxyForHubApiserver(config *rest.Config) (*httputil.ReverseProxy, *url.URL, error) {
	hubApiserverUrl, err := url.Parse(config.Host)
	if err != nil {
		return nil, nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(hubApiserverUrl)
	transport, err := rest.TransportFor(config)
	if err != nil {
		return reverseProxy, hubApiserverUrl, err
	}

	reverseProxy.Transport = transport
	return reverseProxy, hubApiserverUrl, nil
}
