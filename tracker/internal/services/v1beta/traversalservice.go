package v1beta

import (
	"context"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/api/v1beta/graphstore"

	"github.com/golang/protobuf/ptypes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RegisterTraversalServiceServer(server *grpc.Server, graphStore graphstore.GraphStoreClient) {
	v1beta.RegisterTraversalServiceServer(server, &traversalService{
		gs: graphStore,
	})
}

type traversalService struct {
	gs graphstore.GraphStoreClient
}

func (t *traversalService) GetDependents(ctx context.Context, dependency *v1beta.Dependency) (*v1beta.DependentsResponse, error) {
	node, err := newNode(dependency.Module)
	if err != nil {
		return nil, err
	}

	resp, err := t.gs.Neighbors(ctx, &graphstore.NeighborsRequest{
		To: node,
	})
	if err != nil {
		return nil, err
	}

	dependents := make([]*v1beta.Dependency, 0, len(resp.GetNeighbors()))
	for _, neighbor := range resp.GetNeighbors() {
		module := &v1beta.Module{}
		err := ptypes.UnmarshalAny(neighbor.GetNode().GetBody(), module)
		if err != nil {
			continue
		}

		edgeData := make([]*v1beta.ModuleDependency, 0, len(neighbor.GetEdges()))
		for _, edge := range neighbor.GetEdges() {
			moduleDependency := &v1beta.ModuleDependency{}
			err := ptypes.UnmarshalAny(edge.GetBody(), moduleDependency)
			if err != nil {
				continue
			}
			edgeData = append(edgeData, moduleDependency)
		}

		dependents = append(dependents, &v1beta.Dependency{
			Module:   module,
			EdgeData: edgeData,
		})
	}

	return &v1beta.DependentsResponse{
		Dependents: dependents,
	}, nil
}

func (t *traversalService) GetDependencies(ctx context.Context, dependency *v1beta.Dependency) (*v1beta.DependenciesResponse, error) {
	node, err := newNode(dependency.Module)
	if err != nil {
		return nil, err
	}

	resp, err := t.gs.Neighbors(ctx, &graphstore.NeighborsRequest{
		From: node,
	})
	if err != nil {
		return nil, err
	}

	dependencies := make([]*v1beta.Dependency, 0, len(resp.GetNeighbors()))
	for _, neighbor := range resp.GetNeighbors() {
		module := &v1beta.Module{}
		err := ptypes.UnmarshalAny(neighbor.GetNode().GetBody(), module)
		if err != nil {
			continue
		}

		edgeData := make([]*v1beta.ModuleDependency, 0, len(neighbor.GetEdges()))
		for _, edge := range neighbor.GetEdges() {
			moduleDependency := &v1beta.ModuleDependency{}
			err := ptypes.UnmarshalAny(edge.GetBody(), moduleDependency)
			if err != nil {
				continue
			}
			edgeData = append(edgeData, moduleDependency)
		}

		dependencies = append(dependencies, &v1beta.Dependency{
			Module:   module,
			EdgeData: edgeData,
		})
	}

	return &v1beta.DependenciesResponse{
		Dependencies: dependencies,
	}, nil
}

func (t *traversalService) Search(server v1beta.TraversalService_SearchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

func (t *traversalService) BreadthFirstSearch(server v1beta.TraversalService_BreadthFirstSearchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

func (t *traversalService) DepthFirstSearch(server v1beta.TraversalService_DepthFirstSearchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

var _ v1beta.TraversalServiceServer = &traversalService{}
