package service

import (
	dtsapi "github.com/mjpitz/dts/api"
	"golang.org/x/net/context"
)

var _ dtsapi.DependencyTrackingServiceServer = &dependencyTrackingService{}

type dependencyTrackingService struct {

}

func (d *dependencyTrackingService) Put(ctx context.Context, req *dtsapi.PutRequest) (*dtsapi.PutResponse, error) {
	panic("implement me")
}

func (d *dependencyTrackingService) GetDependencies(req *dtsapi.Request, resp dtsapi.DependencyTrackingService_GetDependenciesServer) error {
	panic("implement me")
}

func (d *dependencyTrackingService) GetTopology(req *dtsapi.Request, resp dtsapi.DependencyTrackingService_GetTopologyServer) error {
	panic("implement me")
}

func (d *dependencyTrackingService) GetTopologyTiered(req *dtsapi.Request, resp dtsapi.DependencyTrackingService_GetTopologyTieredServer) error {
	panic("implement me")
}

func (d *dependencyTrackingService) GetSources(req *dtsapi.GetSourcesRequest, resp dtsapi.DependencyTrackingService_GetSourcesServer) error {
	panic("implement me")
}

func (d *dependencyTrackingService) ListLanguages(ctx context.Context, req *dtsapi.ListLanguagesRequest) (*dtsapi.ListLanguagesResponse, error) {
	panic("implement me")
}

func (d *dependencyTrackingService) ListOrganizations(ctx context.Context, req *dtsapi.ListOrganizationsRequest) (*dtsapi.ListOrganizationsResponse, error) {
	panic("implement me")
}

func (d *dependencyTrackingService) ListModules(ctx context.Context, req *dtsapi.ListModulesRequest) (*dtsapi.ListModulesResponse, error) {
	panic("implement me")
}
