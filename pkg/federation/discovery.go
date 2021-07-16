package federation

import (
	"errors"
	"sync"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/discovery"
	cachedDiscoveryClient "k8s.io/client-go/discovery/cached/memory"
)

var (
	apiResourceCache sync.Map
	notFoundError = errors.New("ApiResource not found")
)

func DiscoveryApiResource(info *request.RequestInfo, discoveryClient discovery.DiscoveryInterface) (v1.APIResource, error) {
	groupVersion := schema.GroupVersion{Group: info.APIGroup, Version: info.APIVersion}
	groupVersionKey := groupVersion.String()
	if apiResource, ok := apiResourceCache.Load(groupVersionKey); ok {
		return apiResource.(v1.APIResource), nil
	}

	if discoveryClient == nil {
		return v1.APIResource{}, errors.New("DiscoveryClient is nil")
	}

	cachedClient := cachedDiscoveryClient.NewMemCacheClient(discoveryClient)
	apiResourceList, err := cachedClient.ServerResourcesForGroupVersion(groupVersion.String())
	if err != nil {
		return v1.APIResource{}, err
	}

	if len(apiResourceList.APIResources) <= 0 {
		return v1.APIResource{}, notFoundError
	}

	for _, apiResource := range apiResourceList.APIResources {
		if apiResource.Name == info.Resource {
			cachedApiResource := apiResource.DeepCopy()
			cachedApiResource.Group = info.APIGroup
			cachedApiResource.Version = info.APIVersion
			apiResourceCache.Store(groupVersionKey, *cachedApiResource)
			break
		}
	}

	// get apiResource again, if not exist, return not found error
	if apiResource, ok := apiResourceCache.Load(groupVersionKey); ok {
		return apiResource.(v1.APIResource), nil
	} else {
		return v1.APIResource{}, notFoundError
	}
}
