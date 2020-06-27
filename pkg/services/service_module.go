package services

import (
	"context"

	"github.com/depscloud/api"
	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/tracker/pkg/types"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

// RegisterModuleService registers the moduleService implementation with the server
func RegisterModuleService(server *grpc.Server, gs store.GraphStoreClient) {
	tracker.RegisterModuleServiceServer(server, &moduleService{gs: gs})
}

type moduleService struct {
	gs store.GraphStoreClient
}

var _ tracker.ModuleServiceServer = &moduleService{}

func (s *moduleService) List(ctx context.Context, req *tracker.ListRequest) (*tracker.ListModuleResponse, error) {
	resp, err := s.gs.List(ctx, &store.ListRequest{
		Page:  req.GetPage(),
		Count: req.GetCount(),
		Type:  types.ModuleType,
	})

	if err != nil {
		logrus.Errorf("[service.module] %s", err.Error())
		return nil, err
	}

	modules := make([]*schema.Module, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		module, _ := Decode(item)
		modules = append(modules, module.(*schema.Module))
	}

	return &tracker.ListModuleResponse{
		Page:    req.GetPage(),
		Count:   req.GetCount(),
		Modules: modules,
	}, nil
}

func (s *moduleService) ListSources(ctx context.Context, req *schema.Module) (*tracker.ListSourcesResponse, error) {
	key := keyForModule(req)

	response, err := s.gs.FindDownstream(ctx, &store.FindRequest{
		Key:       key,
		EdgeTypes: []string{types.ManagesType},
	})

	if err != nil {
		return nil, api.ErrModuleNotFound
	}

	sources := make([]*tracker.ManagedSource, len(response.GetPairs()))
	for i, pair := range response.GetPairs() {
		a, _ := Decode(pair.Node)
		b, _ := Decode(pair.Edge)

		sources[i] = &tracker.ManagedSource{
			Source:  a.(*schema.Source),
			Manages: b.(*schema.Manages),
		}
	}

	return &tracker.ListSourcesResponse{
		Sources: sources,
	}, nil
}

func (s *moduleService) ListManaged(ctx context.Context, req *schema.Source) (*tracker.ListManagedResponse, error) {
	key := keyForSource(req)

	response, err := s.gs.FindUpstream(ctx, &store.FindRequest{
		Key:       key,
		EdgeTypes: []string{types.ManagesType},
	})

	if err != nil {
		return nil, api.ErrModuleNotFound
	}

	modules := make([]*tracker.ManagedModule, len(response.GetPairs()))
	for i, pair := range response.GetPairs() {
		a, _ := Decode(pair.Node)
		b, _ := Decode(pair.Edge)

		modules[i] = &tracker.ManagedModule{
			Module:  a.(*schema.Module),
			Manages: b.(*schema.Manages),
		}
	}

	return &tracker.ListManagedResponse{
		Modules: modules,
	}, nil
}
