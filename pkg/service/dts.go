package service

import (
	"fmt"
	dtsapi "github.com/mjpitz/dts/api"
	"github.com/mjpitz/dts/pkg/store"
	"github.com/mjpitz/dts/pkg/types"
	"golang.org/x/net/context"
)

var _ dtsapi.DependencyTrackingServiceServer = &dependencyTrackingService{}

type dependencyTrackingService struct {
	graphStore store.GraphStore
}

func quickKey(gi *store.GraphItem) string {
	return fmt.Sprintf("%s:%s:%s",
		gi.GraphItemType, string(gi.K1), string(gi.K2))
}

func (d *dependencyTrackingService) Put(ctx context.Context, req *dtsapi.PutRequest) (*dtsapi.PutResponse, error) {
	graphItems := types.ExtractGraphItems(req)

	currentIndex := make(map[string]*store.GraphItem)
	for _, gi := range graphItems {
		currentIndex[quickKey(gi)] = gi
	}

	sourceGraphItem := graphItems[0]
	managedModules, err := d.graphStore.FindDownstream(sourceGraphItem.K1)
	if err != nil {
		return nil, err
	}

	toRemove := make([]*store.PrimaryKey, 0)

	for _, managedModule := range managedModules {
		_, managedExists := currentIndex[quickKey(managedModule)]

		if !managedExists {
			toRemove = append(toRemove, &store.PrimaryKey{
				GraphItemType: types.ManagesType,
				K1:            sourceGraphItem.K1,
				K2:            managedModule.K1,
			})

			continue
		}

		dependedModules, err := d.graphStore.FindDownstream(managedModule.K1)
		if err != nil {
			return nil, err
		}

		for _, dependedModule := range dependedModules {
			_, dependedExists := currentIndex[quickKey(dependedModule)]

			if !dependedExists {
				toRemove = append(toRemove, &store.PrimaryKey{
					GraphItemType: types.DependsType,
					K1:            managedModule.K1,
					K2:            dependedModule.K2,
				})
			}
		}
	}

	if err := d.graphStore.Delete(toRemove); err != nil {
		return nil, err
	}

	if err := d.graphStore.Put(graphItems); err != nil {
		return nil, err
	}

	return &dtsapi.PutResponse{
		Code:    "dts-200",
		Message: "Success",
	}, nil
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
