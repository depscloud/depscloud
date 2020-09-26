package services

import (
	"context"

	"github.com/depscloud/api"
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/tracker/internal/graphstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterSearchService registers the searchService implementation with the server
func RegisterSearchService(server *grpc.Server, gs store.GraphStoreClient) {
	tracker.RegisterSearchServiceServer(server, &searchService{
		gs: gs,
		ss: &sourceService{gs: gs},
		ms: &moduleService{gs: gs},
		ds: &dependencyService{gs: gs},
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
	stream := make(chan *tracker.SearchRequest, 2)

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

	return nil
}

func transformResponse(response *tracker.SearchResponse) (requests []*tracker.SearchRequest, err error) {
	if dependencies := response.GetDependencies(); dependencies != nil {
		requests = make([]*tracker.SearchRequest, len(dependencies))

		for i, dependency := range dependencies {
			requests[i] = &tracker.SearchRequest{
				DependenciesOf: &tracker.DependencyRequest{
					Language:     dependency.Module.GetLanguage(),
					Organization: dependency.Module.GetOrganization(),
					Module:       dependency.Module.GetModule(),
				},
			}
		}

	} else if dependents := response.GetDependents(); dependents != nil {
		requests = make([]*tracker.SearchRequest, len(dependents))

		for i, dependency := range dependents {
			requests[i] = &tracker.SearchRequest{
				DependentsOf: &tracker.DependencyRequest{
					Language:     dependency.Module.GetLanguage(),
					Organization: dependency.Module.GetOrganization(),
					Module:       dependency.Module.GetModule(),
				},
			}
		}
	} else {
		err = api.ErrUnimplemented
	}

	return requests, err
}

func keyFor(request *tracker.SearchRequest) string {
	if dependenciesOf := request.GetDependenciesOf(); dependenciesOf != nil {
		return graphstore.Base64encode(keyForDependencyRequest(dependenciesOf))

	} else if dependentsOf := request.GetDependentsOf(); dependentsOf != nil {
		return graphstore.Base64encode(keyForDependencyRequest(dependentsOf))

	} else if sourcesFor := request.GetSourcesFor(); sourcesFor != nil {
		return graphstore.Base64encode(keyForModule(sourcesFor))

	} else if modulesFor := request.GetModulesFor(); modulesFor != nil {
		return graphstore.Base64encode(keyForSource(modulesFor))
	}

	return ""
}

func (s *searchService) BreadthFirstSearch(server tracker.SearchService_BreadthFirstSearchServer) (err error) {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	done := ctx.Done()
	stream := consumeStream(ctx, server)

	root := <-stream
	queue := []*tracker.SearchRequest{
		root,
	}
	seen := map[string]bool{
		keyFor(root): true,
	}

	for length := len(queue); length > 0; length = len(queue) {
		next := make([]*tracker.SearchRequest, 0)
		toSend := make([]*tracker.SearchResponse, length)

		// process tier
		for i := 0; i < length; i++ {
			request := queue[i]

			toSend[i], err = s.processRequest(ctx, request)
			if err != nil {
				return err
			}

			nextBatch, err := transformResponse(toSend[i])
			if err != nil {
				return err
			}

			for _, item := range nextBatch {
				key := keyFor(item)
				if _, ok := seen[key]; !ok {
					seen[key] = true
					next = append(next, item)
				}
			}
		}

		// flush responses
		for _, response := range toSend {
			if err := server.Send(response); err != nil {
				return err
			}
		}

		select {
		case <-done:
			return nil
		case req := <-stream:
			if req.GetCancel() {
				return nil
			}

			return status.Error(codes.InvalidArgument, "unexpected request body")
		default:
			queue = next
		}
	}

	return nil
}

func (s *searchService) DepthFirstSearch(server tracker.SearchService_DepthFirstSearchServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	done := ctx.Done()
	stream := consumeStream(ctx, server)

	root := <-stream
	stack := []*tracker.SearchRequest{
		root,
	}
	seen := map[string]bool{
		keyFor(root): true,
	}

	for length := len(stack); length > 0; length = len(stack) {
		next := make([]*tracker.SearchRequest, 0)

		// Pop
		node := stack[length-1]
		stack = stack[0:length]

		// Explore Node
		var response *tracker.SearchResponse
		response, err := s.processRequest(ctx, node)
		if err != nil {
			return err
		}

		// Get adjacent nodes
		nextBatch, err := transformResponse(response)
		if err != nil {
			return err
		}

		for _, item := range nextBatch {
			key := keyFor(item)
			if _, ok := seen[key]; !ok {
				seen[key] = true
				next = append(next, item)
			}
		}

		// send response to client
		if err := server.Send(response); err != nil {
			return err
		}

		select {
		case <-done:
			return nil

		// TODO: What this case is for?
		// Is it to handle some unexpected request from cmd line when the original request is being processed?
		case req := <-stream:
			if req.GetCancel() {
				return nil
			}

			return status.Error(codes.InvalidArgument, "unexpected request body")
		default:
			stack = next
		}
	}

	return nil
}

var _ tracker.SearchServiceServer = &searchService{}
