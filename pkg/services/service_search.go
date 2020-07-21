package services

import (
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/api/v1alpha/tracker"

	"google.golang.org/grpc"
)

// RegisterSearchService registers the searchService implementation with the server
func RegisterSearchService(server *grpc.Server, gs store.GraphStoreClient) {
	tracker.RegisterSearchServiceServer(server, &searchService{gs: gs})
}

type searchService struct {
	tracker.UnimplementedSearchServiceServer

	gs store.GraphStoreClient
}

var _ tracker.SearchServiceServer = &searchService{}
