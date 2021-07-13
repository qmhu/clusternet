package main

import (
	"bytes"
	"fmt"
	"github.com/clusternet/clusternet/pkg/utils"
	"io/ioutil"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var (
	openAddr = "127.0.0.1:2002"
	configPath = "/Users/hu/.kube/config"
	apiserverUrl = "https://127.0.0.1:54869"
)


func main() {
	apiserverURI, _ := url.Parse(apiserverUrl)

	childKubeConfig, _ := utils.LoadsKubeConfig(configPath, 1)
    var transport http.RoundTripper
	transport, _ = rest.TransportFor(childKubeConfig)
	proxy := httputil.NewSingleHostReverseProxy(apiserverURI)
	proxy.Transport = transport

	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", apiserverUrl)

		req.URL.Scheme = apiserverURI.Scheme
		req.URL.Host = apiserverURI.Host

		klog.Infof("Proxy request %s", formatRequest(req))
	}

	proxy.ModifyResponse = func(r *http.Response) error {
		bodyInRes := r.Body
		defer bodyInRes.Close()
		bodyBytes, err := ioutil.ReadAll(bodyInRes)
		if err != nil {
			return err
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		klog.Infof("Proxy response: %v, status %s", string(bodyBytes), r.Status)
		return nil
	}

	log.Println("Server is starting ：" + openAddr)
	log.Fatalln(http.ListenAndServe(openAddr, proxy))
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
