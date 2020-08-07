package proxies

import (
	"context"
	"io"

	"github.com/depscloud/api/v1alpha/tracker"
)

func NewSearchServiceProxy(client tracker.SearchServiceClient) tracker.SearchServiceServer {
	return &searchService{
		client: client,
	}
}

type searchClientStream interface {
	Send(*tracker.SearchRequest) error
	Recv() (*tracker.SearchResponse, error)
}

type searchServerStream interface {
	Send(*tracker.SearchResponse) error
	Recv() (*tracker.SearchRequest, error)
}

func process(parent context.Context, client searchClientStream, server searchServerStream) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	done := ctx.Done()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				req, err := server.Recv()
				if err != nil {
					return
				}

				if err = client.Send(req); err != nil {
					return
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			return nil
		default:
			resp, err := client.Recv()
			if err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}

			if err = server.Send(resp); err != nil {
				return err
			}
		}
	}
}

type searchService struct {
	client tracker.SearchServiceClient
}

func (s *searchService) Search(server tracker.SearchService_SearchServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	call, err := s.client.Search(ctx)
	if err != nil {
		return err
	}

	return process(ctx, call, server)
}

func (s *searchService) BreadthFirstSearch(server tracker.SearchService_BreadthFirstSearchServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	call, err := s.client.BreadthFirstSearch(ctx)
	if err != nil {
		return err
	}

	return process(ctx, call, server)
}

func (s *searchService) DepthFirstSearch(server tracker.SearchService_DepthFirstSearchServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	call, err := s.client.DepthFirstSearch(ctx)
	if err != nil {
		return err
	}

	return process(ctx, call, server)
}

var _ tracker.SearchServiceServer = &searchService{}
