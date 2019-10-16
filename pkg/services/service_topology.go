package services

import (
	"github.com/deps-cloud/api"
	"github.com/deps-cloud/api/v1alpha/store"
	"github.com/deps-cloud/api/v1alpha/tracker"

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

func (t *topologyService) GetDependentsTopology(req *tracker.DependencyRequest, resp tracker.TopologyService_GetDependentsTopologyServer) error {
	return api.ErrUnimplemented
}

func (t *topologyService) GetDependentsTopologyTiered(req *tracker.DependencyRequest, resp tracker.TopologyService_GetDependentsTopologyTieredServer) error {
	return api.ErrUnimplemented
}

func (t *topologyService) GetDependenciesTopology(req *tracker.DependencyRequest, resp tracker.TopologyService_GetDependenciesTopologyServer) error {
	return api.ErrUnimplemented
}

func (t *topologyService) GetDependenciesTopologyTiered(req *tracker.DependencyRequest, resp tracker.TopologyService_GetDependenciesTopologyTieredServer) error {
	return api.ErrUnimplemented
}
