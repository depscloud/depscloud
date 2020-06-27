package services

import (
	"context"

	"github.com/depscloud/api"
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/api/v1alpha/tracker"

	"google.golang.org/grpc"
)

// RegisterTopologyService registers the topologyService implementation with the server
func RegisterTopologyService(server *grpc.Server, gs store.GraphStoreClient) {
	tracker.RegisterTopologyServiceServer(server, &topologyService{gs: gs})
}

type topologyService struct {
	gs store.GraphStoreClient
}

var _ tracker.TopologyServiceServer = &topologyService{}

func (t *topologyService) ListDependentsTopology(ctx context.Context, req *tracker.DependencyRequest) (*tracker.ListDependentsResponse, error) {
	return nil, api.ErrUnimplemented
}

func (t *topologyService) ListDependentsTopologyTiered(ctx context.Context, req *tracker.DependencyRequest) (*tracker.ListDependentsTieredResponse, error) {
	return nil, api.ErrUnimplemented
}

func (t *topologyService) ListDependenciesTopology(ctx context.Context, req *tracker.DependencyRequest) (*tracker.ListDependenciesResponse, error) {
	return nil, api.ErrUnimplemented
}

func (t *topologyService) ListDependenciesTopologyTiered(ctx context.Context, req *tracker.DependencyRequest) (*tracker.ListDependenciesTieredResponse, error) {
	return nil, api.ErrUnimplemented
}
