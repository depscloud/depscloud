package service

import (
	"fmt"
	"net/http"

	dtsapi "github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/pkg/store"
	"github.com/deps-cloud/dts/pkg/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// NewDependencyTrackingService constructs a service using the specified graphstore.
func NewDependencyTrackingService(graphStore store.GraphStore) (dtsapi.DependencyTrackerServer, error) {
	return &dependencyTrackingService{
		graphStore: graphStore,
	}, nil
}

var _ dtsapi.DependencyTrackerServer = &dependencyTrackingService{}

type dependencyTrackingService struct {
	graphStore store.GraphStore
}

func quickKey(gi *store.GraphItem) string {
	return fmt.Sprintf("%s:%s:%s",
		gi.GraphItemType, string(gi.K1), string(gi.K2))
}

func (d *dependencyTrackingService) Put(ctx context.Context, req *dtsapi.PutRequest) (*dtsapi.PutResponse, error) {
	url := req.GetSourceInformation().GetUrl()

	traversalUtil := &TraversalUtil{ d.graphStore, dtsapi.Direction_DOWNSTREAM }
	graphItems := types.ExtractGraphItems(req)

	currentIndex := make(map[string]*store.GraphItem)
	for _, gi := range graphItems {
		currentIndex[quickKey(gi)] = gi
	}

	sourceGraphItem := graphItems[0]
	managedModules, err := traversalUtil.GetAdjacent(sourceGraphItem.K1, []string{ types.ManagesType })
	if err != nil {
		logrus.Errorf("failed to fetch managed modules: %s, %v", url, err)
		return nil, dtsapi.ErrModuleNotFound
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

		dependedModules, err := traversalUtil.GetAdjacent(managedModule.K1, []string{ types.DependsType })
		if err != nil {
			logrus.Errorf("failed to fetch depended modules: %s, %v", url, err)
			return nil, dtsapi.ErrModuleNotFound
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
		logrus.Errorf("[service.dts] failed to delete removed edges: %v", err)
		return nil, dtsapi.ErrPartialDeletion
	}

	if err := d.graphStore.Put(graphItems); err != nil {
		logrus.Errorf("[service.dts] failed to add new edges: %v", err)
		return nil, dtsapi.ErrPartialInsertion
	}

	return &dtsapi.PutResponse{
		Code:    http.StatusOK,
		Message: "successfully updated source",
	}, nil
}

func (d *dependencyTrackingService) GetDependencies(req *dtsapi.Request, resp dtsapi.DependencyTracker_GetDependenciesServer) error {
	url := fmt.Sprintf("%s://%s;%s", req.Language, req.Organization, req.Module)
	logrus.Infof("looking up dependencies for %s", url)

	traversalUtil := &TraversalUtil{ d.graphStore, req.Direction }
	key := types.ExtractModuleKey(req)

	dependencies, err := traversalUtil.GetAdjacent(key, []string{ types.DependsType })
	if err != nil {
		logrus.Errorf("failed to fetch dependencies: %s, %v", url, err)
		return dtsapi.ErrModuleNotFound
	}

	for _, dep := range dependencies {
		item, err := types.Decode(dep)
		if err != nil {
			// type / encoding problem, skip
			logrus.Errorf("[service.dts] failed to decode dependency: %v", err)
			continue
		}

		module := item.(*types.Module)
		response := &dtsapi.Response{
			Dependency: &dtsapi.DependencyId{
				Language: module.Language,
				Organization: module.Organization,
				Module: module.Module,
			},
		}

		if err = resp.Send(response); err != nil {
			logrus.Errorf("[service.dts] failed to send response: %v", err)
		}
	}

	return nil
}

func (d *dependencyTrackingService) GetManaged(ctx context.Context, req *dtsapi.GetManagedRequest) (*dtsapi.GetManagedResponse, error) {
	traversalUtil := &TraversalUtil{ d.graphStore, dtsapi.Direction_DOWNSTREAM }
	key := types.ExtractSourceKey(req)

	managed, err := traversalUtil.GetAdjacent(key, []string{ types.ManagesType })
	if err != nil {
		logrus.Errorf("failed to fetch managed: %s, %v", req.Url, err)
		return nil, dtsapi.ErrModuleNotFound
	}

	depIds := make([]*dtsapi.DependencyId, 0, len(managed))
	for _, dep := range managed {
		item, err := types.Decode(dep)

		if err != nil {
			// type / encoding problem, skip
			logrus.Errorf("[service.dts] failed to decode dependency: %v", err)
			continue
		}

		module := item.(*types.Module)
		depIds = append(depIds, &dtsapi.DependencyId{
			Language: module.Language,
			Organization: module.Organization,
			Module: module.Module,
		})
	}

	return &dtsapi.GetManagedResponse{
		Url: req.Url,
		Managed: depIds,
	}, nil
}

func (d *dependencyTrackingService) GetTopology(req *dtsapi.Request, resp dtsapi.DependencyTracker_GetTopologyServer) error {
	return dtsapi.ErrUnimplemented
}

func (d *dependencyTrackingService) GetTopologyTiered(req *dtsapi.Request, resp dtsapi.DependencyTracker_GetTopologyTieredServer) error {
	return dtsapi.ErrUnimplemented
}

func (d *dependencyTrackingService) GetSources(req *dtsapi.GetSourcesRequest, resp dtsapi.DependencyTracker_GetSourcesServer) error {
	return fmt.Errorf("unimplemented")
}

func (d *dependencyTrackingService) ListLanguages(ctx context.Context, req *dtsapi.ListLanguagesRequest) (*dtsapi.ListLanguagesResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (d *dependencyTrackingService) ListOrganizations(ctx context.Context, req *dtsapi.ListOrganizationsRequest) (*dtsapi.ListOrganizationsResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (d *dependencyTrackingService) ListModules(ctx context.Context, req *dtsapi.ListModulesRequest) (*dtsapi.ListModulesResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}
