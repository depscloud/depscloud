package v1beta

import (
	"context"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/api/v1beta/graphstore"
	"github.com/depscloud/depscloud/internal/logger"

	"github.com/golang/protobuf/ptypes"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

func RegisterModuleServiceServer(server *grpc.Server, graphStore graphstore.GraphStoreClient) {
	v1beta.RegisterModuleServiceServer(server, &moduleService{
		gs: graphStore,
	})
}

type moduleService struct {
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
		return nil, err
	}

	modules := make([]*v1beta.Module, 0, len(resp.GetNodes()))
	for _, node := range resp.GetNodes() {
		module := &v1beta.Module{}
		err := ptypes.UnmarshalAny(node.GetBody(), module)
		if err != nil {
			log.Error("failed to parse module", zap.Error(err))
			continue
		}
		modules = append(modules, module)
	}

	return &v1beta.ListModulesResponse{
		NextPageToken: resp.GetNextPageToken(),
		Modules:       modules,
	}, nil
}

func (m *moduleService) ListSources(ctx context.Context, module *v1beta.ManagedModule) (*v1beta.ListManagedSourcesResponse, error) {
	node, err := newNode(module.Module)
	if err != nil {
		return nil, err
	}

	resp, err := m.gs.Neighbors(ctx, &graphstore.NeighborsRequest{
		To: node,
	})
	if err != nil {
		return nil, err
	}

	sources := make([]*v1beta.ManagedSource, 0, len(resp.GetNeighbors()))
	for _, neighbor := range resp.GetNeighbors() {
		source := &v1beta.Source{}
		err := ptypes.UnmarshalAny(neighbor.GetNode().GetBody(), source)
		if err != nil {
			continue
		}

		edgeData := make([]*v1beta.SourceModule, 0, len(neighbor.GetEdges()))
		for _, edge := range neighbor.GetEdges() {
			sourceModule := &v1beta.SourceModule{}
			err := ptypes.UnmarshalAny(edge.GetBody(), sourceModule)
			if err != nil {
				continue
			}
			edgeData = append(edgeData, sourceModule)
		}

		sources = append(sources, &v1beta.ManagedSource{
			Source:   source,
			EdgeData: edgeData,
		})
	}

	return &v1beta.ListManagedSourcesResponse{
		Sources: sources,
	}, nil
}

var _ v1beta.ModuleServiceServer = &moduleService{}
