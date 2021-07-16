package handler

import (
	"context"
	"encoding/json"
	"github.com/clusternet/clusternet/pkg/federation/handler/manifestgetter"
	"github.com/clusternet/clusternet/pkg/federation/handler/responser"
	"io/ioutil"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	registryrest "k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/rest"

	clusternetClientSet "github.com/clusternet/clusternet/pkg/generated/clientset/versioned"
)

type UpdateHandler struct {
	apiResource      metav1.APIResource
	requestInfo      *request.RequestInfo
	clusternetClient *clusternetClientSet.Clientset
	config           *rest.Config
}

func NewUpdateHandler(apiResource metav1.APIResource, requestInfo *request.RequestInfo, clusternetClient *clusternetClientSet.Clientset, config *rest.Config) *UpdateHandler {
	return &UpdateHandler{
		apiResource:      apiResource,
		requestInfo:      requestInfo,
		clusternetClient: clusternetClient,
		config:           config,
	}
}

func (h *UpdateHandler) HandleConnection(ctx context.Context, w http.ResponseWriter, req *http.Request, responder registryrest.Responder) error {
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return apierrors.NewBadRequest(err.Error())
	}

	requestBodyObj := &unstructured.Unstructured{}
	err = json.Unmarshal(requestBody, requestBodyObj)
	if err != nil {
		return apierrors.NewBadRequest(err.Error())
	}

	manifest, err := manifestgetter.GetManifest(ctx, h.apiResource, h.requestInfo, h.clusternetClient)
	if err != nil {
		return err
	}

	manifest.Template = runtime.RawExtension{Raw: requestBody}


	updated, err := h.clusternetClient.AppsV1alpha1().Manifests(manifest.Namespace).Update(ctx, manifest, metav1.UpdateOptions{})
	if err != nil {
		return apierrors.NewInternalError(err)
	}

	feedObj, err := responser.ManifestToObject(updated)
	if err != nil {
		return apierrors.NewInternalError(err)
	}

	responser.WriteObjectNegotiated(h.apiResource, feedObj, w, req, responder)
	return nil
}
