package services

import (
	"github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/api/v1alpha"
	"github.com/deps-cloud/dts/api/v1alpha/store"

	"google.golang.org/grpc"
)

// RegisterTopologyService registers the topologyService implementation with the server
func RegisterTopologyService(server *grpc.Server, gs store.GraphStoreClient) {
	v1alpha.RegisterTopologyServiceServer(server, &topologyService{gs: gs})
}

type topologyService struct {
	gs store.GraphStoreClient
}

var _ v1alpha.TopologyServiceServer = &topologyService{}

func (t *topologyService) GetDependentsTopology(req *v1alpha.DependencyRequest, resp v1alpha.TopologyService_GetDependentsTopologyServer) error {
	return api.ErrUnimplemented
}

func (t *topologyService) GetDependentsTopologyTiered(req *v1alpha.DependencyRequest, resp v1alpha.TopologyService_GetDependentsTopologyTieredServer) error {
	return api.ErrUnimplemented
}

func (t *topologyService) GetDependenciesTopology(req *v1alpha.DependencyRequest, resp v1alpha.TopologyService_GetDependenciesTopologyServer) error {
	return api.ErrUnimplemented
}

func (t *topologyService) GetDependenciesTopologyTiered(req *v1alpha.DependencyRequest, resp v1alpha.TopologyService_GetDependenciesTopologyTieredServer) error {
	return api.ErrUnimplemented
}
