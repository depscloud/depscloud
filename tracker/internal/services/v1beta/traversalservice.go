package v1beta

import (
	"context"
	"fmt"

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
	ms v1beta.ModuleServiceClient
	ss v1beta.SourceServiceClient
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
	ctx := server.Context()
	call, err := t.gs.Traverse(ctx)
	if err != nil {
		return err
	}

	for {
		var err error

		req, err := server.Recv()
		if err != nil {
			return err
		}

		if req.Cancel {
			return call.Send(&graphstore.TraverseRequest{
				Cancel: true,
			})
		}

		resp := &v1beta.SearchResponse{
			Request: req,
		}

		if req.DependenciesFor != nil {
			r, e := t.GetDependencies(ctx, req.DependenciesFor)

			resp.Dependencies = r.GetDependencies()
			err = e
		} else if req.DependentsOf != nil {
			r, e := t.GetDependents(ctx, req.DependentsOf)

			resp.Dependents = r.GetDependents()
			err = e
		} else if req.ModulesFor != nil {
			r, e := t.ss.ListModules(ctx, req.ModulesFor)

			resp.Modules = r.GetModules()
			err = e
		} else if req.SourcesOf != nil {
			r, e := t.ms.ListSources(ctx, req.SourcesOf)

			resp.Sources = r.GetSources()
			err = e
		} else {
			return fmt.Errorf("unrecognized request")
		}

		if err != nil {
			return err
		}

		err = server.Send(resp)
		if err != nil {
			return err
		}
	}
}

func (t *traversalService) BreadthFirstSearch(server v1beta.TraversalService_BreadthFirstSearchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

func (t *traversalService) DepthFirstSearch(server v1beta.TraversalService_DepthFirstSearchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

var _ v1beta.TraversalServiceServer = &traversalService{}
