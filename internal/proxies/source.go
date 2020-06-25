package proxies

import (
	"context"

	"github.com/deps-cloud/api/v1alpha/tracker"
)

func NewSourceServiceProxy(client tracker.SourceServiceClient) tracker.SourceServiceServer {
	return &sourceService{
		client: client,
	}
}

type sourceService struct {
	client tracker.SourceServiceClient
}

func (s *sourceService) List(ctx context.Context, request *tracker.ListRequest) (*tracker.ListSourceResponse, error) {
	return s.client.List(ctx, request)
}

func (s *sourceService) Track(ctx context.Context, request *tracker.SourceRequest) (*tracker.TrackResponse, error) {
	return s.client.Track(ctx, request)
}

var _ tracker.SourceServiceServer = &sourceService{}
