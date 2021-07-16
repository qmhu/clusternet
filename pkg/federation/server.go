package federation

import (
	"context"
	"net/http"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/request"
	registryrest "k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	federationsv1alpha1 "github.com/clusternet/clusternet/pkg/apis/federations/v1alpha1"
	"github.com/clusternet/clusternet/pkg/federation/constants"
	"github.com/clusternet/clusternet/pkg/federation/handler"
	clusternetClientSet "github.com/clusternet/clusternet/pkg/generated/clientset/versioned"
)

type Server struct {
	config           *rest.Config
	kubeClient       *kubernetes.Clientset
	clusternetClient *clusternetClientSet.Clientset
	ApiserverConfig  *genericapiserver.RecommendedConfig
}

func NewServer(config *rest.Config, kubeClient *kubernetes.Clientset, clusternetClient *clusternetClientSet.Clientset) *Server {
	return &Server{
		config:           config,
		kubeClient:       kubeClient,
		clusternetClient: clusternetClient,
	}
}

func (s *Server) HandleConnection(ctx context.Context, id string, opts *federationsv1alpha1.Declaration, responder registryrest.Responder) (http.Handler, error) {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// get declarations urls
		declarationsIndex := strings.Index(request.URL.Path, constants.RequestPathDeclarations)
		if declarationsIndex < 0 {
			responder.Error(errors.NewBadRequest("Cannot get declarations path from request"))
		}

		declarationsUrl := request.URL.Path[declarationsIndex+len(constants.RequestPathDeclarations):]
		declarationsReq := request.Clone(ctx)
		declarationsReq.URL.Path = declarationsUrl

		// retrieve request info based on declarations url
		requestInfo, err := s.ApiserverConfig.RequestInfoResolver.NewRequestInfo(declarationsReq)
		if err != nil {
			responder.Error(errors.NewBadRequest("Cannot get request info from request"))
			return
		}

		// get api resource for declarations
		apiResource := v1.APIResource{}
		if requestInfo.IsResourceRequest {
			apiResource, err = DiscoveryApiResource(requestInfo, s.kubeClient.DiscoveryClient)
			if err != nil {
				responder.Error(err)
				return
			}
		}

		requestHandler := s.RouteRequestHandler(requestInfo, apiResource, declarationsUrl)
		err = requestHandler.HandleConnection(ctx, writer, request, responder)
		if err != nil {
			klog.Errorf("handle request %s error %v", request.URL.String(), err)
			responder.Error(err)
		}
	}), nil
}

func (s *Server) RouteRequestHandler(requestInfo *request.RequestInfo, apiResource v1.APIResource, declarationsUrl string) handler.RequestHandler {
	if requestInfo.IsResourceRequest {
		switch strings.ToLower(requestInfo.Verb) {
		case "get":
			return handler.NewGetHandler(apiResource, requestInfo, s.clusternetClient)
		case "create":
			return handler.NewCreationHandler(apiResource, requestInfo, s.clusternetClient, s.config, declarationsUrl)
		case "update":
			return handler.NewUpdateHandler(apiResource, requestInfo, s.clusternetClient, s.config)
		case "patch":
			return handler.NewPatchHandler(apiResource, requestInfo, s.clusternetClient, s.config)
		case "list":
			return handler.NewListHandler(apiResource, requestInfo, s.clusternetClient)
		default:
			return handler.NewDummyHandler()
		}
	} else {
		// forward non-resource requests to hub apiserver
		return handler.NewProxyHandler(s.config, declarationsUrl)
	}
}
