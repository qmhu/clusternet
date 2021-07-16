package manifestgetter

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apiserver/pkg/endpoints/request"

	"github.com/clusternet/clusternet/pkg/apis/apps/v1alpha1"
	"github.com/clusternet/clusternet/pkg/federation/constants"
	clusternetClientSet "github.com/clusternet/clusternet/pkg/generated/clientset/versioned"
)

func GetManifestNamespace(resource v1.APIResource, requestNamespace string) string {
	if resource.Namespaced {
		return requestNamespace
	} else {
		return constants.ClusternetSystemNamespace
	}

}

func GetManifest(ctx context.Context, resource v1.APIResource, info *request.RequestInfo, clientset *clusternetClientSet.Clientset) (*v1alpha1.Manifest, error) {
	namespace := GetManifestNamespace(resource, info.Namespace)

	labelSelector := labels.NewSelector()
	feedSelectorList := make(map[string]string)
	feedSelectorList[constants.FeedKindLabel] = resource.Kind
	feedSelectorList[constants.FeedApiVersionLabel] = resource.Version
	feedSelectorList[constants.FeedNameLabel] = info.Name
	feedSelectorList[constants.FeedNamespaceLabel] = namespace
	for key, value := range feedSelectorList {
		requirement, err := labels.NewRequirement(key, selection.Equals, []string{value})
		if err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}

		labelSelector = labelSelector.Add(*requirement)
	}

	manifestList, err := clientset.AppsV1alpha1().Manifests(namespace).List(ctx, v1.ListOptions{
		LabelSelector: labelSelector.String(),
	})

	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	if len(manifestList.Items) <= 0 {
		return nil, errors.NewBadRequest("Cannot found manifest")
	}
	return &manifestList.Items[0], nil
}
