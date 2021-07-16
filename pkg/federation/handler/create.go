package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/endpoints/request"
	registryrest "k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	appsapi "github.com/clusternet/clusternet/pkg/apis/apps/v1alpha1"
	"github.com/clusternet/clusternet/pkg/federation/constants"
	"github.com/clusternet/clusternet/pkg/federation/handler/manifestgetter"
	clusternetClientSet "github.com/clusternet/clusternet/pkg/generated/clientset/versioned"
	"github.com/clusternet/clusternet/pkg/known"
)

type CreationHandler struct {
	apiResource      metav1.APIResource
	requestInfo      *request.RequestInfo
	clusternetClient *clusternetClientSet.Clientset
	config           *rest.Config
	proxyUrl         string
}

func NewCreationHandler(apiResource metav1.APIResource, requestInfo *request.RequestInfo, clusternetClient *clusternetClientSet.Clientset, config *rest.Config, proxyUrl string) *CreationHandler {
	return &CreationHandler{
		apiResource:      apiResource,
		requestInfo:      requestInfo,
		clusternetClient: clusternetClient,
		config:           config,
		proxyUrl:         proxyUrl,
	}
}

func (h *CreationHandler) HandleConnection(ctx context.Context, w http.ResponseWriter, req *http.Request, responder registryrest.Responder) error {
	reverseProxy, hubApiserverUrl, err := NewReverseProxyForHubApiserver(h.config)
	if err != nil {
		return apierrors.NewInternalError(err)
	}

	requestBodyObj := &unstructured.Unstructured{}
	var requestBody []byte

	reverseProxy.Director = func(req *http.Request) {
		req.URL.Path = h.proxyUrl
		req.URL.Scheme = hubApiserverUrl.Scheme
		req.URL.Host = hubApiserverUrl.Host

		klog.V(4).Infof("Add DryRun options to validation object, uri %s", req.URL.String())
		DryRun(req)

		// read request body and save it
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			klog.Error("Cannot read body from request")
			return
		}

		requestBody = body
		err = json.Unmarshal(requestBody, requestBodyObj)
		if err != nil {
			klog.Error("Unmarshal request body failed")
		}

		// Write back body
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	reverseProxy.ModifyResponse = func(res *http.Response) error {
		bodyInRes := res.Body
		defer bodyInRes.Close()
		bodyBytes, err := ioutil.ReadAll(bodyInRes)
		if err != nil {
			return err
		}

		// Request has pass the dryRun validation
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			err = wait.PollImmediate(constants.ManifestRequestInterval, constants.ManifestRequestTimeout, func() (bool, error) {
				// Name generation
				manifestName := fmt.Sprintf("%s-%s", strings.ToLower(requestBodyObj.GetKind()), rand.String(constants.ManifestInstanceIDLength))
				namespace := manifestgetter.GetManifestNamespace(h.apiResource, h.requestInfo.Namespace)
				manifest := &appsapi.Manifest{
					ObjectMeta: metav1.ObjectMeta{
						Name:      manifestName,
						Namespace: namespace,
						Labels: map[string]string{
							known.ObjectCreatedByLabel:    known.ClusternetHubName,
							constants.FeedApiVersionLabel: requestBodyObj.GetAPIVersion(),
							constants.FeedKindLabel:       requestBodyObj.GetKind(),
							constants.FeedNameLabel:       requestBodyObj.GetName(),
							constants.FeedNamespaceLabel:  namespace,
						},
						Annotations: requestBodyObj.GetAnnotations(),
						Finalizers: []string{
							known.AppFinalizer,
						},
					},
					Template: runtime.RawExtension{Raw: requestBody},
				}

				for key, value := range requestBodyObj.GetLabels() {
					manifest.Labels[key] = value
				}

				_, err := h.clusternetClient.AppsV1alpha1().Manifests(namespace).Create(ctx, manifest, metav1.CreateOptions{})
				if err != nil {
					if apierrors.IsAlreadyExists(err) {
						return false, nil
					} else {
						return false, err
					}
				}

				return true, nil
			})

			if err != nil {
				return apierrors.NewInternalError(err)
			}
		}

		res.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		klog.V(4).Infof("Handle federation, response: %v, status %s", string(bodyBytes), res.Status)

		return nil
	}

	reverseProxy.ErrorHandler = proxy.NewErrorResponder(responder).Error
	reverseProxy.ServeHTTP(w, req)
	return nil
}
