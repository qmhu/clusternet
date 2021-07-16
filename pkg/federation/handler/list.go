package handler

import (
	"context"
	"fmt"
	"github.com/clusternet/clusternet/pkg/federation/handler/manifestgetter"
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apiserver/pkg/endpoints/request"
	registryrest "k8s.io/apiserver/pkg/registry/rest"

	"github.com/clusternet/clusternet/pkg/federation/constants"
	"github.com/clusternet/clusternet/pkg/federation/handler/responser"
	clusternetClientSet "github.com/clusternet/clusternet/pkg/generated/clientset/versioned"
)

type ListHandler struct {
	apiResource      metav1.APIResource
	requestInfo      *request.RequestInfo
	clusternetClient *clusternetClientSet.Clientset
}

func NewListHandler(apiResource metav1.APIResource, requestInfo *request.RequestInfo, clusternetClient *clusternetClientSet.Clientset) *ListHandler {
	return &ListHandler{
		apiResource:      apiResource,
		requestInfo:      requestInfo,
		clusternetClient: clusternetClient,
	}
}

func (h *ListHandler) HandleConnection(ctx context.Context, w http.ResponseWriter, req *http.Request, responder registryrest.Responder) error {
	gv := schema.GroupVersion{Group: h.apiResource.Group, Version: h.apiResource.Version}
	requestScope := responser.GetRequestScope(responder)

	opts := metav1.ListOptions{}

	if err := requestScope.ParameterCodec.DecodeParameters(req.URL.Query(), gv, &opts); err != nil {
		return errors.NewBadRequest(err.Error())
	}

	mergedLabelSelector, err := labels.Parse(opts.LabelSelector)
	if err != nil {
		return errors.NewBadRequest(err.Error())
	}

	namespace := manifestgetter.GetManifestNamespace(h.apiResource, h.requestInfo.Namespace)

	feedSelectorList := make(map[string]string)
	feedSelectorList[constants.FeedKindLabel] = h.apiResource.Kind
	feedSelectorList[constants.FeedApiVersionLabel] = h.apiResource.Version
	feedSelectorList[constants.FeedNamespaceLabel] = namespace
	for key, value := range feedSelectorList {
		requirement, err := labels.NewRequirement(key, selection.Equals, []string{value})
		if err != nil {
			return errors.NewBadRequest(err.Error())
		}

		mergedLabelSelector = mergedLabelSelector.Add(*requirement)
	}

	manifestList, err := h.clusternetClient.AppsV1alpha1().Manifests(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: mergedLabelSelector.String(),
	})

	if err != nil {
		return errors.NewInternalError(err)
	}

	unstructedList := &unstructured.UnstructuredList{}
	for _, manifest := range manifestList.Items {
		feedObj, _ := responser.ManifestToObject(&manifest)
		item := feedObj.(*unstructured.Unstructured)

		unstructedList.Items = append(unstructedList.Items, *item)
	}

	unstructedList.SetAPIVersion(h.apiResource.Version)
	unstructedList.SetKind(fmt.Sprintf("%sList", h.apiResource.Kind))

	responser.WriteObjectNegotiated(h.apiResource, unstructedList, w, req, responder)
	return nil
}
