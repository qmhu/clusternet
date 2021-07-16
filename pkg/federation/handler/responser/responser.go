package responser

import (
	"encoding/json"
	"net/http"
	"reflect"
	"unsafe"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	registryrest "k8s.io/apiserver/pkg/registry/rest"

	"github.com/clusternet/clusternet/pkg/apis/apps/v1alpha1"
)

func WriteObjectNegotiated(apiResource v1.APIResource, obj runtime.Object, w http.ResponseWriter, req *http.Request, responder registryrest.Responder) {
	gv := schema.GroupVersion{Group: apiResource.Group, Version: apiResource.Version}
	requestScope := GetRequestScope(responder)
	responsewriters.WriteObjectNegotiated(requestScope.Serializer, requestScope, gv, w, req, 200, obj)
}

func ManifestToObject(manifest *v1alpha1.Manifest) (runtime.Object, error) {
	feedObj := &unstructured.Unstructured{}
	err := json.Unmarshal(manifest.Template.Raw, feedObj)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	feedObj.SetCreationTimestamp(manifest.GetCreationTimestamp())
	feedObj.SetResourceVersion(manifest.GetResourceVersion())
	feedObj.SetDeletionTimestamp(manifest.GetDeletionTimestamp())
	feedObj.SetUID(manifest.GetUID())

	return feedObj, nil
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func GetRequestScope(responder registryrest.Responder) *handlers.RequestScope {
	responderValue := reflect.ValueOf(responder)
	return GetUnexportedField(responderValue.Elem().FieldByName("scope")).(*handlers.RequestScope)
}
