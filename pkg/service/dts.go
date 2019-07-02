package service

import (
	"fmt"
	"net/http"

	dtsapi "github.com/deps-cloud/tracker/api"
	"github.com/deps-cloud/tracker/api/v1alpha/schema"
	"github.com/deps-cloud/tracker/api/v1alpha/store"
	"github.com/deps-cloud/tracker/pkg/services"
	"github.com/deps-cloud/tracker/pkg/types"

	"github.com/sirupsen/logrus"

	"golang.org/x/net/context"
)

// NewDependencyTrackingService constructs a service using the specified graphstore.
func NewDependencyTrackingService(graphStore store.GraphStoreClient) (dtsapi.DependencyTrackerServer, error) {
	return &dependencyTrackingService{
		graphStore: graphStore,
	}, nil
}

var _ dtsapi.DependencyTrackerServer = &dependencyTrackingService{}

type dependencyTrackingService struct {
	graphStore store.GraphStoreClient
}

func quickKey(gi *store.GraphItem) string {
	return fmt.Sprintf("%s:%s:%s",
		gi.GraphItemType, string(gi.K1), string(gi.K2))
}

func (d *dependencyTrackingService) Put(ctx context.Context, req *dtsapi.PutRequest) (*dtsapi.PutResponse, error) {
	url := req.GetSourceInformation().GetUrl()

	traversalUtil := &TraversalUtil{d.graphStore, dtsapi.Direction_UPSTREAM}
	graphItems := ExtractGraphItems(req)

	currentIndex := make(map[string]*store.GraphItem)
	for _, gi := range graphItems {
		currentIndex[quickKey(gi)] = gi
	}

	sourceGraphItem := graphItems[0]
	managedModules, err := traversalUtil.GetAdjacent(sourceGraphItem.GetK1(), []string{types.ManagesType})
	if err != nil {
		logrus.Errorf("failed to fetch managed modules: %s, %v", url, err)
		return nil, dtsapi.ErrModuleNotFound
	}

	toRemove := make([]*store.GraphItem, 0)

	for _, managedModule := range managedModules {
		_, managedExists := currentIndex[quickKey(managedModule)]

		if !managedExists {
			toRemove = append(toRemove, managedModule)

			continue
		}

		dependedModules, err := traversalUtil.GetAdjacent(managedModule.GetK1(), []string{types.DependsType})
		if err != nil {
			logrus.Errorf("failed to fetch depended modules: %s, %v", url, err)
			return nil, dtsapi.ErrModuleNotFound
		}

		for _, dependedModule := range dependedModules {
			_, dependedExists := currentIndex[quickKey(dependedModule)]

			if !dependedExists {
				toRemove = append(toRemove, dependedModule)
			}
		}
	}

	if _, err := d.graphStore.Delete(ctx, &store.DeleteRequest{Items: toRemove}); err != nil {
		logrus.Errorf("[service.source] %s", err.Error())
		return nil, dtsapi.ErrPartialDeletion
	}

	if _, err := d.graphStore.Put(ctx, &store.PutRequest{Items: graphItems}); err != nil {
		logrus.Errorf("[service.source] %s", err.Error())
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

	traversalUtil := &TraversalUtil{d.graphStore, req.Direction}
	key := ExtractModuleKeyFromRequest(req)

	dependencies, err := traversalUtil.GetAdjacent(key, []string{types.DependsType})
	if err != nil {
		logrus.Errorf("failed to fetch dependencies: %s, %v", url, err)
		return dtsapi.ErrModuleNotFound
	}

	for _, dep := range dependencies {
		item, err := services.Decode(dep)
		if err != nil {
			// type / encoding problem, skip
			logrus.Errorf("[service.dts] failed to decode dependency: %v", err)
			continue
		}

		module := item.(*schema.Module)
		response := &dtsapi.Response{
			Dependency: &dtsapi.DependencyId{
				Language:     module.GetLanguage(),
				Organization: module.GetOrganization(),
				Module:       module.GetModule(),
			},
		}

		if err = resp.Send(response); err != nil {
			logrus.Errorf("[service.dts] failed to send response: %v", err)
		}
	}

	return nil
}

func (d *dependencyTrackingService) GetManaged(ctx context.Context, req *dtsapi.GetManagedRequest) (*dtsapi.GetManagedResponse, error) {
	key := ExtractSourceKey(req)

	managed, err := d.graphStore.FindUpstream(context.Background(), &store.FindRequest{
		Key:       key,
		EdgeTypes: []string{types.ManagesType},
	})
	if err != nil {
		logrus.Errorf("failed to fetch managed: %s, %v", req.Url, err)
		return nil, dtsapi.ErrModuleNotFound
	}

	depIds := make([]*dtsapi.DependencyId, 0, len(managed.GetPairs()))
	for _, dep := range managed.GetPairs() {
		item, err := services.Decode(dep.GetNode())

		if err != nil {
			// type / encoding problem, skip
			logrus.Errorf("[service.dts] failed to decode dependency: %v", err)
			continue
		}

		module := item.(*schema.Module)
		depIds = append(depIds, &dtsapi.DependencyId{
			Language:     module.GetLanguage(),
			Organization: module.GetOrganization(),
			Module:       module.GetModule(),
		})
	}

	return &dtsapi.GetManagedResponse{
		Url:     req.Url,
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
	key := ExtractModuleKeyFromGetSourcesRequest(req)

	sources, err := d.graphStore.FindDownstream(context.Background(), &store.FindRequest{
		Key:       key,
		EdgeTypes: []string{types.ManagesType},
	})
	if err != nil {
		logrus.Errorf("failed to fetch sources: %s, %v", req, err)
		return dtsapi.ErrModuleNotFound
	}

	for _, source := range sources.GetPairs() {
		item, err := services.Decode(source.GetNode())
		if err != nil {
			// type / encoding problem, skip
			logrus.Errorf("[service.dts] failed to decode source: %v", err)
			continue
		}

		source := item.(*schema.Source)
		response := &dtsapi.GetSourcesResponse{
			Source: &dtsapi.SourceInformation{
				Url: source.GetUrl(),
			},
		}

		if err = resp.Send(response); err != nil {
			logrus.Errorf("[service.dts] failed to send response: %v", err)
		}
	}

	return nil
}

func (d *dependencyTrackingService) ListLanguages(ctx context.Context, req *dtsapi.ListLanguagesRequest) (*dtsapi.ListLanguagesResponse, error) {
	return nil, dtsapi.ErrUnimplemented
}

func (d *dependencyTrackingService) ListOrganizations(ctx context.Context, req *dtsapi.ListOrganizationsRequest) (*dtsapi.ListOrganizationsResponse, error) {
	return nil, dtsapi.ErrUnimplemented
}

func (d *dependencyTrackingService) ListModules(ctx context.Context, req *dtsapi.ListModulesRequest) (*dtsapi.ListModulesResponse, error) {
	return nil, dtsapi.ErrUnimplemented
}
