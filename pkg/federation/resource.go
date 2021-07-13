package federation

import (
	"k8s.io/apiserver/pkg/endpoints/request"
	"net/http"
	"strings"
)

func IsResourceRequest(requestInfo *request.RequestInfo) bool {
	return requestInfo.IsResourceRequest
}

func IsUpdateMethod(requestInfo *request.RequestInfo) bool {
	if strings.ToLower(requestInfo.Verb) == "create" ||
		strings.ToLower(requestInfo.Verb) == "update" ||
		strings.ToLower(requestInfo.Verb) == "patch" {
		return true
	}

	return false
}

func DryRun(req *http.Request) {
	q := req.URL.Query()
	q.Add("dryRun", "All")
	req.URL.RawQuery = q.Encode()
}