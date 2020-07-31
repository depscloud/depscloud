package services

import (
	"context"
	"github.com/depscloud/api"
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/api/v1alpha/tracker"
	"google.golang.org/grpc"
)

// RegisterSearchService registers the searchService implementation with the server
func RegisterSearchService(server *grpc.Server, gs store.GraphStoreClient) {
	tracker.RegisterSearchServiceServer(server, &searchService{
		gs: gs,
	})
}

type searchService struct {
	gs store.GraphStoreClient

	ss tracker.SourceServiceServer
	ms tracker.ModuleServiceServer
	ds tracker.DependencyServiceServer
}

type consumable interface {
	Recv() (*tracker.SearchRequest, error)
}

func consumeStream(ctx context.Context, c consumable) chan *tracker.SearchRequest {
	done := ctx.Done()
	stream := make(chan *tracker.SearchRequest)

	go func() {
		select {
		case <-done:
			break
		default:
			req, err := c.Recv()
			if err != nil {
				break
			}

			stream <- req
		}
	}()

	return stream
}

func (s *searchService) processRequest(ctx context.Context, request *tracker.SearchRequest) (response *tracker.SearchResponse, interr error) {
	response = &tracker.SearchResponse{
		Request: request,
	}

	if dependenciesOf := request.GetDependenciesOf(); dependenciesOf != nil {
		result, err := s.ds.ListDependencies(ctx, dependenciesOf)
		interr = err

		response.Dependencies = result.GetDependencies()

	} else if dependentsOf := request.GetDependentsOf(); dependentsOf != nil {
		result, err := s.ds.ListDependents(ctx, dependentsOf)
		interr = err

		response.Dependents = result.GetDependents()

	} else if sourcesFor := request.GetSourcesFor(); sourcesFor != nil {
		result, err := s.ms.ListSources(ctx, sourcesFor)
		interr = err

		response.Sources = result.GetSources()

	} else if modulesFor := request.GetModulesFor(); modulesFor != nil {
		result, err := s.ms.ListManaged(ctx, modulesFor)
		interr = err

		response.Modules = result.GetModules()

	} else {
		interr = api.ErrUnimplemented
	}

	return response, interr
}

func (s *searchService) Search(server tracker.SearchService_SearchServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	done := ctx.Done()
	stream := consumeStream(ctx, server)

	select {
	case <-done:
		return nil
	case request := <-stream:
		if request.GetCancel() {
			return nil
		}

		// TODO: eventually send this to a task queue

		response, err := s.processRequest(ctx, request)
		if err != nil {
			return err
		}

		if err := server.Send(response); err != nil {
			return err
		}
	}

	return api.ErrUnimplemented
}

func (s *searchService) BreadthFirstSearch(server tracker.SearchService_BreadthFirstSearchServer) error {
	return api.ErrUnimplemented
}

func (s *searchService) DepthFirstSearch(server tracker.SearchService_DepthFirstSearchServer) error {
	return api.ErrUnimplemented
}

var _ tracker.SearchServiceServer = &searchService{}
