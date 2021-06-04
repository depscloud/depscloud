package v1beta

import (
	"context"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/api/v1beta/graphstore"
	"github.com/depscloud/depscloud/internal/logger"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

func RegisterModuleServiceServer(server *grpc.Server, graphStore graphstore.GraphStoreClient) {
	v1beta.RegisterModuleServiceServer(server, &moduleService{
		gs: graphStore,
	})
}

type moduleService struct {
	v1beta.UnsafeModuleServiceServer

	gs graphstore.GraphStoreClient
}

func (m *moduleService) List(ctx context.Context, request *v1beta.ListRequest) (*v1beta.ListModulesResponse, error) {
	log := logger.Extract(ctx)

	resp, err := m.gs.List(ctx, &graphstore.ListRequest{
		Parent:    request.GetParent(),
		PageSize:  request.GetPageSize(),
		PageToken: request.GetPageToken(),
		Kind:      moduleKind,
	})
	if err != nil {
		log.Error("failed to list modules", zap.Error(err))
		return nil, ErrInvalidRequest
	}

	modules := make([]*v1beta.Module, 0, len(resp.GetNodes()))
	for _, node := range resp.GetNodes() {
		module, err := fromNodeOrEdge(node, &v1beta.Module{})
		if err != nil {
			log.Warn("failed to parse module", zap.Error(err))
			continue
		}
		modules = append(modules, module.(*v1beta.Module))
	}

	return &v1beta.ListModulesResponse{
		NextPageToken: resp.GetNextPageToken(),
		Modules:       modules,
	}, nil
}

func (m *moduleService) ListSources(ctx context.Context, module *v1beta.ManagedModule) (*v1beta.ListManagedSourcesResponse, error) {
	log := logger.Extract(ctx)

	node, err := newNode(module.Module)
	if err != nil {
		log.Error("failed to parse module into node", zap.Error(err))
		return nil, ErrInvalidRequest
	}

	resp, err := m.gs.Neighbors(ctx, &graphstore.NeighborsRequest{
		To: node,
	})
	if err != nil {
		log.Error("failed to query graph", zap.Error(err))
		return nil, ErrQueryFailure
	}

	sources := make([]*v1beta.ManagedSource, 0, len(resp.GetNeighbors()))
	for _, neighbor := range resp.GetNeighbors() {
		managedSource, errors := neighborToManagedSource(neighbor)

		for _, err := range errors {
			log.Warn("encountered an issue converting managed source", zap.Error(err))
		}

		if managedSource != nil {
			sources = append(sources, managedSource)
		}
	}

	return &v1beta.ListManagedSourcesResponse{
		Sources: sources,
	}, nil
}

var _ v1beta.ModuleServiceServer = &moduleService{}

func neighborToManagedSource(neighbor *graphstore.Neighbor) (_ *v1beta.ManagedSource, errors []error) {
	source, err := fromNodeOrEdge(neighbor.GetNode(), &v1beta.Source{})
	if err != nil {
		return nil, []error{err}
	}

	edgeData := make([]*v1beta.SourceModule, 0, len(neighbor.GetEdges()))
	for _, edge := range neighbor.GetEdges() {
		sourceModule, err := fromNodeOrEdge(edge, &v1beta.SourceModule{})
		if err != nil {
			errors = append(errors, err)
		} else {
			edgeData = append(edgeData, sourceModule.(*v1beta.SourceModule))
		}
	}

	return &v1beta.ManagedSource{
		Source:   source.(*v1beta.Source),
		EdgeData: edgeData,
	}, errors
}
