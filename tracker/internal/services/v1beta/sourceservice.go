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

func RegisterSourceServiceServer(server *grpc.Server, graphStore graphstore.GraphStoreClient) {
	v1beta.RegisterSourceServiceServer(server, &sourceService{
		gs: graphStore,
	})
}

type sourceService struct {
	gs graphstore.GraphStoreClient
}

func (s *sourceService) List(ctx context.Context, request *v1beta.ListRequest) (*v1beta.ListSourcesResponse, error) {
	log := logger.Extract(ctx)

	resp, err := s.gs.List(ctx, &graphstore.ListRequest{
		Parent:    request.GetParent(),
		PageSize:  request.GetPageSize(),
		PageToken: request.GetPageToken(),
		Kind:      sourceKind,
	})
	if err != nil {
		return nil, err
	}

	sources := make([]*v1beta.Source, 0, len(resp.GetNodes()))
	for _, node := range resp.GetNodes() {
		source := &v1beta.Source{}
		err := ptypes.UnmarshalAny(node.GetBody(), source)
		if err != nil {
			log.Error("failed parse source", zap.Error(err))
			continue
		}
		sources = append(sources, source)
	}

	return &v1beta.ListSourcesResponse{
		NextPageToken: resp.GetNextPageToken(),
		Sources:       sources,
	}, nil
}

func (s *sourceService) ListModules(ctx context.Context, source *v1beta.ManagedSource) (*v1beta.ListManagedModulesResponse, error) {
	node, err := newNode(source.Source)
	if err != nil {
		return nil, err
	}

	resp, err := s.gs.Neighbors(ctx, &graphstore.NeighborsRequest{
		From: node,
	})
	if err != nil {
		return nil, err
	}

	modules := make([]*v1beta.ManagedModule, 0, len(resp.GetNeighbors()))
	for _, neighbor := range resp.GetNeighbors() {
		module := &v1beta.Module{}
		err := ptypes.UnmarshalAny(neighbor.GetNode().GetBody(), module)
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

		modules = append(modules, &v1beta.ManagedModule{
			Module:   module,
			EdgeData: edgeData,
		})
	}

	return &v1beta.ListManagedModulesResponse{
		Modules: modules,
	}, nil
}

var _ v1beta.SourceServiceServer = &sourceService{}
