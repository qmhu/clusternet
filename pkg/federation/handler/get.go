package handler

import (
	"context"
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/request"
	registryrest "k8s.io/apiserver/pkg/registry/rest"

	"github.com/clusternet/clusternet/pkg/federation/handler/manifestgetter"
	"github.com/clusternet/clusternet/pkg/federation/handler/responser"
	clusternetClientSet "github.com/clusternet/clusternet/pkg/generated/clientset/versioned"
)

type GetHandler struct {
	apiResource      metav1.APIResource
	requestInfo      *request.RequestInfo
	clusternetClient *clusternetClientSet.Clientset
}

func NewGetHandler(apiResource metav1.APIResource, requestInfo *request.RequestInfo, clusternetClient *clusternetClientSet.Clientset) *GetHandler {
	return &GetHandler{
		apiResource:      apiResource,
		requestInfo:      requestInfo,
		clusternetClient: clusternetClient,
	}
}

func (h *GetHandler) HandleConnection(ctx context.Context, w http.ResponseWriter, req *http.Request, responder registryrest.Responder) error {
	manifest, err := manifestgetter.GetManifest(ctx, h.apiResource, h.requestInfo, h.clusternetClient)
	if err != nil {
		return err
	}

	feedObj, err := responser.ManifestToObject(manifest)
	if err != nil {
		return errors.NewInternalError(err)
	}

	responser.WriteObjectNegotiated(h.apiResource, feedObj, w, req, responder)
	return nil
}
