package federation

import (
	"bytes"
	"context"
	"fmt"
	types "github.com/clusternet/clusternet/pkg/apis/federations/v1alpha1"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apiserver/pkg/endpoints/request"
	registryrest "k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ReverseProxyHandler struct {
	config *rest.Config
	transport http.RoundTripper
	hubApiserverUrl *url.URL
	ctx context.Context
	id string
	opts *types.Govern
	responder registryrest.Responder
}

func NewReverseProxyHandler(config *rest.Config, ctx context.Context, id string, opts *types.Govern, responder registryrest.Responder) *ReverseProxyHandler {
	hubApiserverUrl, _ := url.Parse(config.Host)

	transport, _ := rest.TransportFor(config)

	return &ReverseProxyHandler{
		config: config,
		hubApiserverUrl: hubApiserverUrl,
		transport: transport,
		ctx: ctx,
		id: id,
		opts: opts,
		responder: responder,
	}
}

func (r *ReverseProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reverseProxy := httputil.NewSingleHostReverseProxy(r.hubApiserverUrl)
	requestInfo, _ := request.RequestInfoFrom(r.ctx)

	reverseProxy.Transport = r.transport
	isDryRun := false

	reverseProxy.Director = func(req *http.Request) {
		req.Header.Add("X-Clusternet-Forwarded-Host", req.Host)
		req.Header.Add("X-Clusternet-Origin-Host", r.hubApiserverUrl.Host)

		klog.Infof("Handle federation,origin uri %s", req.URL.String())

		governIndex := strings.Index(req.URL.Path, "governs")
		proxyUrl := req.URL.Path[governIndex+7:]

		req.URL.Path = proxyUrl
		req.URL.Scheme = r.hubApiserverUrl.Scheme
		req.URL.Host = r.hubApiserverUrl.Host

		if IsResourceRequest(requestInfo) && IsUpdateMethod(requestInfo) {
			klog.Infof("Add DryRun options to validation object, uri %s", req.URL.String())
			DryRun(req)
			isDryRun = true
		}

		klog.Infof("Handle federation, requesturi %s, requestinfo %v", formatRequest(req), requestInfo)
	}

	reverseProxy.ModifyResponse = func(r *http.Response) error {
		bodyInRes := r.Body
		defer bodyInRes.Close()
		bodyBytes, err := ioutil.ReadAll(bodyInRes)
		if err != nil {
			return err
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		klog.Infof("Handle federation, response: %v, status %s", string(bodyBytes), r.Status)

		if isDryRun {

		}

		return nil
	}
	reverseProxy.ErrorHandler = proxy.NewErrorResponder(r.responder).Error

	reverseProxy.ServeHTTP(w, req)
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
