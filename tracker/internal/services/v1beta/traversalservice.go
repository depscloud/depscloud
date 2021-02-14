package v1beta

import (
	"context"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/api/v1beta/graphstore"
	"github.com/depscloud/depscloud/internal/logger"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

func RegisterTraversalServiceServer(server *grpc.Server, graphStore graphstore.GraphStoreClient) {
	v1beta.RegisterTraversalServiceServer(server, &traversalService{
		gs: graphStore,
		ms: &moduleService{gs: graphStore},
		ss: &sourceService{gs: graphStore},
	})
}

type traversalService struct {
	gs graphstore.GraphStoreClient
	ms v1beta.ModuleServiceServer
	ss v1beta.SourceServiceServer
}

func (t *traversalService) GetDependents(ctx context.Context, dependency *v1beta.Dependency) (*v1beta.DependentsResponse, error) {
	log := logger.Extract(ctx)

	node, err := newNode(dependency.Module)
	if err != nil {
		log.Error("failed to convert module to node", zap.Error(err))
		return nil, ErrInvalidRequest
	}

	resp, err := t.gs.Neighbors(ctx, &graphstore.NeighborsRequest{
		To: node,
	})
	if err != nil {
		log.Error("failed to retrieve dependents", zap.Error(err))
		return nil, ErrQueryFailure
	}

	dependents := make([]*v1beta.Dependency, 0, len(resp.GetNeighbors()))
	for _, neighbor := range resp.GetNeighbors() {
		dependency, errors := neighborToDependency(neighbor)

		for _, err := range errors {
			log.Warn("encountered an issue converting dependent", zap.Error(err))
		}

		if dependency != nil {
			dependents = append(dependents, dependency)
		}
	}

	return &v1beta.DependentsResponse{
		Dependents: dependents,
	}, nil
}

func (t *traversalService) GetDependencies(ctx context.Context, dependency *v1beta.Dependency) (*v1beta.DependenciesResponse, error) {
	log := logger.Extract(ctx)

	node, err := newNode(dependency.Module)
	if err != nil {
		log.Error("failed to convert module to node", zap.Error(err))
		return nil, ErrInvalidRequest
	}

	resp, err := t.gs.Neighbors(ctx, &graphstore.NeighborsRequest{
		From: node,
	})
	if err != nil {
		log.Error("failed to retrieve dependencies", zap.Error(err))
		return nil, ErrQueryFailure
	}

	dependencies := make([]*v1beta.Dependency, 0, len(resp.GetNeighbors()))
	for _, neighbor := range resp.GetNeighbors() {
		dependency, errors := neighborToDependency(neighbor)

		for _, err := range errors {
			log.Warn("encountered an issue converting dependency", zap.Error(err))
		}

		if dependency != nil {
			dependencies = append(dependencies, dependency)
		}
	}

	return &v1beta.DependenciesResponse{
		Dependencies: dependencies,
	}, nil
}

var _ v1beta.TraversalServiceServer = &traversalService{}

// neighborToDependency is a helper function used to convert a neighbor structure to a dependency structure.
func neighborToDependency(neighbor *graphstore.Neighbor) (*v1beta.Dependency, []error) {
	module, err := fromNodeOrEdge(neighbor.GetNode(), &v1beta.Module{})
	if err != nil {
		return nil, []error{err}
	}

	errors := make([]error, 0)

	edgeData := make([]*v1beta.ModuleDependency, 0, len(neighbor.GetEdges()))
	for _, edge := range neighbor.GetEdges() {
		moduleDependency, err := fromNodeOrEdge(edge, &v1beta.ModuleDependency{})
		if err != nil {
			errors = append(errors, err)
		} else {
			edgeData = append(edgeData, moduleDependency.(*v1beta.ModuleDependency))
		}
	}

	return &v1beta.Dependency{
		Module:   module.(*v1beta.Module),
		EdgeData: edgeData,
	}, nil
}
