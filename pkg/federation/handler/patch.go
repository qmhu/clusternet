package handler

import (
	"context"
	"fmt"
	"github.com/clusternet/clusternet/pkg/federation/handler/manifestgetter"
	"github.com/clusternet/clusternet/pkg/federation/handler/responser"
	"io/ioutil"
	"net/http"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/endpoints/request"
	registryrest "k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/rest"

	clusternetClientSet "github.com/clusternet/clusternet/pkg/generated/clientset/versioned"
)

type PatchHandler struct {
	apiResource      metav1.APIResource
	requestInfo      *request.RequestInfo
	clusternetClient *clusternetClientSet.Clientset
	config           *rest.Config
}

func NewPatchHandler(apiResource metav1.APIResource, requestInfo *request.RequestInfo, clusternetClient *clusternetClientSet.Clientset, config *rest.Config) *PatchHandler {
	return &PatchHandler{
		apiResource:      apiResource,
		requestInfo:      requestInfo,
		clusternetClient: clusternetClient,
		config:           config,
	}
}

func (h *PatchHandler) HandleConnection(ctx context.Context, w http.ResponseWriter, req *http.Request, responder registryrest.Responder) error {
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return apierrors.NewBadRequest(err.Error())
	}

	// TODO: handle this in negotiation
	contentType := req.Header.Get("Content-Type")
	// Remove "; charset=" if included in header.
	if idx := strings.Index(contentType, ";"); idx > 0 {
		contentType = contentType[:idx]
	}
	patchType := types.PatchType(contentType)
	if patchType != types.StrategicMergePatchType {
		return apierrors.NewBadRequest("only support strategic-merge-patch now")
	}

	// only support smp now
	patchBytes := []byte(fmt.Sprintf(`{"template": %s }`, string(requestBody)))

	manifest, err := manifestgetter.GetManifest(ctx, h.apiResource, h.requestInfo, h.clusternetClient)
	if err != nil {
		return err
	}

	patched, err := h.clusternetClient.AppsV1alpha1().Manifests(manifest.Namespace).Patch(ctx, manifest.Name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		return apierrors.NewInternalError(err)
	}

	feedObj, err := responser.ManifestToObject(patched)
	if err != nil {
		return apierrors.NewInternalError(err)
	}

	responser.WriteObjectNegotiated(h.apiResource, feedObj, w, req, responder)
	return nil
}
