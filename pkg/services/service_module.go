package services

import (
	"context"

	"github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/api/v1alpha"
	"github.com/deps-cloud/dts/api/v1alpha/schema"
	"github.com/deps-cloud/dts/api/v1alpha/store"
	"github.com/deps-cloud/dts/pkg/types"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

// RegisterModuleService registers the moduleService implementation with the server
func RegisterModuleService(server *grpc.Server, gs store.GraphStoreClient) {
	v1alpha.RegisterModuleServiceServer(server, &moduleService{gs: gs})
}

type moduleService struct {
	gs store.GraphStoreClient
}

var _ v1alpha.ModuleServiceServer = &moduleService{}

func (s *moduleService) List(ctx context.Context, req *v1alpha.ListRequest) (*v1alpha.ListModuleResponse, error) {
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

	return &v1alpha.ListModuleResponse{
		Page:    req.GetPage(),
		Count:   req.GetCount(),
		Modules: modules,
	}, nil
}

func (s *moduleService) GetSource(req *schema.Module, resp v1alpha.ModuleService_GetSourceServer) error {
	key := keyForModule(req)

	response, err := s.gs.FindDownstream(context.Background(), &store.FindRequest{
		Key:       key,
		EdgeTypes: []string{types.ManagesType},
	})

	if err != nil {
		return api.ErrModuleNotFound
	}

	for _, pair := range response.GetPairs() {
		a, _ := Decode(pair.Node)
		b, _ := Decode(pair.Edge)

		managedSource := &v1alpha.ManagedSource{
			Source:  a.(*schema.Source),
			Manages: b.(*schema.Manages),
		}

		if err := resp.Send(managedSource); err != nil {
			logrus.Errorf("[service.dependency] failed to send response: %v", err)
		}
	}

	return nil
}

func (s *moduleService) GetManaged(req *schema.Source, resp v1alpha.ModuleService_GetManagedServer) error {
	key := keyForSource(req)

	response, err := s.gs.FindUpstream(context.Background(), &store.FindRequest{
		Key:       key,
		EdgeTypes: []string{types.ManagesType},
	})

	if err != nil {
		return api.ErrModuleNotFound
	}

	for _, pair := range response.GetPairs() {
		a, _ := Decode(pair.Node)
		b, _ := Decode(pair.Edge)

		managedModule := &v1alpha.ManagedModule{
			Module:  a.(*schema.Module),
			Manages: b.(*schema.Manages),
		}

		if err := resp.Send(managedModule); err != nil {
			logrus.Errorf("[service.dependency] failed to send response: %v", err)
		}
	}

	return nil
}
