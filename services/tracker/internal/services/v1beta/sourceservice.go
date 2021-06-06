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

func RegisterSourceServiceServer(server *grpc.Server, graphStore graphstore.GraphStoreClient, index IndexService) {
	v1beta.RegisterSourceServiceServer(server, &sourceService{
		gs:    graphStore,
		index: index,
	})
}

type sourceService struct {
	v1beta.UnsafeSourceServiceServer

	gs    graphstore.GraphStoreClient
	index IndexService
}

func (s *sourceService) Search(ctx context.Context, request *v1beta.SourcesSearchRequest) (*v1beta.ListSourcesResponse, error) {
	log := logger.Extract(ctx)

	results, err := s.index.Query(ctx, &Index{
		Kind:  sourceKind,
		Field: "url",
		Value: request.Like.Url,
	})
	if err != nil {
		log.Error("failed to search modules", zap.Error(err))
		return nil, ErrQueryFailure
	}

	sources := make([]*v1beta.Source, 0, len(results))
	for _, result := range results {
		sources = append(sources, &v1beta.Source{
			Url: result.Value,
		})
	}

	return &v1beta.ListSourcesResponse{
		Sources: sources,
	}, nil
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
		log.Error("failed to lookup sources", zap.Error(err))
		return nil, ErrQueryFailure
	}

	sources := make([]*v1beta.Source, 0, len(resp.GetNodes()))
	for _, node := range resp.GetNodes() {
		source := &v1beta.Source{}
		err := ptypes.UnmarshalAny(node.GetBody(), source)
		if err != nil {
			log.Warn("failed parse source", zap.Error(err))
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
	log := logger.Extract(ctx)

	node, err := newNode(source.Source)
	if err != nil {
		log.Error("failed to parse source into node", zap.Error(err))
		return nil, ErrInvalidRequest
	}

	resp, err := s.gs.Neighbors(ctx, &graphstore.NeighborsRequest{
		From: node,
	})
	if err != nil {
		log.Error("failed to query for modules", zap.Error(err))
		return nil, ErrQueryFailure
	}

	modules := make([]*v1beta.ManagedModule, 0, len(resp.GetNeighbors()))
	for _, neighbor := range resp.GetNeighbors() {
		managedModule, errors := neighborToManagedModule(neighbor)

		for _, err := range errors {
			log.Warn("encountered an issue converting managed module", zap.Error(err))
		}

		if managedModule != nil {
			modules = append(modules, managedModule)
		}
	}

	return &v1beta.ListManagedModulesResponse{
		Modules: modules,
	}, nil
}

var _ v1beta.SourceServiceServer = &sourceService{}

func neighborToManagedModule(neighbor *graphstore.Neighbor) (_ *v1beta.ManagedModule, errors []error) {
	module, err := fromNodeOrEdge(neighbor.GetNode(), &v1beta.Module{})
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

	return &v1beta.ManagedModule{
		Module:   module.(*v1beta.Module),
		EdgeData: edgeData,
	}, errors
}
