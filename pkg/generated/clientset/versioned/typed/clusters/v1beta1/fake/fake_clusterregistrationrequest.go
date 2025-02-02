/*
Copyright 2021 The Clusternet Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/clusternet/clusternet/pkg/apis/clusters/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterRegistrationRequests implements ClusterRegistrationRequestInterface
type FakeClusterRegistrationRequests struct {
	Fake *FakeClustersV1beta1
}

var clusterregistrationrequestsResource = schema.GroupVersionResource{Group: "clusters.clusternet.io", Version: "v1beta1", Resource: "clusterregistrationrequests"}

var clusterregistrationrequestsKind = schema.GroupVersionKind{Group: "clusters.clusternet.io", Version: "v1beta1", Kind: "ClusterRegistrationRequest"}

// Get takes name of the clusterRegistrationRequest, and returns the corresponding clusterRegistrationRequest object, and an error if there is any.
func (c *FakeClusterRegistrationRequests) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.ClusterRegistrationRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterregistrationrequestsResource, name), &v1beta1.ClusterRegistrationRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ClusterRegistrationRequest), err
}

// List takes label and field selectors, and returns the list of ClusterRegistrationRequests that match those selectors.
func (c *FakeClusterRegistrationRequests) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.ClusterRegistrationRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterregistrationrequestsResource, clusterregistrationrequestsKind, opts), &v1beta1.ClusterRegistrationRequestList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.ClusterRegistrationRequestList{ListMeta: obj.(*v1beta1.ClusterRegistrationRequestList).ListMeta}
	for _, item := range obj.(*v1beta1.ClusterRegistrationRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterRegistrationRequests.
func (c *FakeClusterRegistrationRequests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusterregistrationrequestsResource, opts))
}

// Create takes the representation of a clusterRegistrationRequest and creates it.  Returns the server's representation of the clusterRegistrationRequest, and an error, if there is any.
func (c *FakeClusterRegistrationRequests) Create(ctx context.Context, clusterRegistrationRequest *v1beta1.ClusterRegistrationRequest, opts v1.CreateOptions) (result *v1beta1.ClusterRegistrationRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusterregistrationrequestsResource, clusterRegistrationRequest), &v1beta1.ClusterRegistrationRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ClusterRegistrationRequest), err
}

// Update takes the representation of a clusterRegistrationRequest and updates it. Returns the server's representation of the clusterRegistrationRequest, and an error, if there is any.
func (c *FakeClusterRegistrationRequests) Update(ctx context.Context, clusterRegistrationRequest *v1beta1.ClusterRegistrationRequest, opts v1.UpdateOptions) (result *v1beta1.ClusterRegistrationRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusterregistrationrequestsResource, clusterRegistrationRequest), &v1beta1.ClusterRegistrationRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ClusterRegistrationRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeClusterRegistrationRequests) UpdateStatus(ctx context.Context, clusterRegistrationRequest *v1beta1.ClusterRegistrationRequest, opts v1.UpdateOptions) (*v1beta1.ClusterRegistrationRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(clusterregistrationrequestsResource, "status", clusterRegistrationRequest), &v1beta1.ClusterRegistrationRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ClusterRegistrationRequest), err
}

// Delete takes name of the clusterRegistrationRequest and deletes it. Returns an error if one occurs.
func (c *FakeClusterRegistrationRequests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(clusterregistrationrequestsResource, name), &v1beta1.ClusterRegistrationRequest{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterRegistrationRequests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusterregistrationrequestsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.ClusterRegistrationRequestList{})
	return err
}

// Patch applies the patch and returns the patched clusterRegistrationRequest.
func (c *FakeClusterRegistrationRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.ClusterRegistrationRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusterregistrationrequestsResource, name, pt, data, subresources...), &v1beta1.ClusterRegistrationRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ClusterRegistrationRequest), err
}
