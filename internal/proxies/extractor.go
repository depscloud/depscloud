package proxies

import (
	"context"

	"github.com/deps-cloud/api/v1alpha/extractor"
)

func NewExtractorServiceProxy(client extractor.DependencyExtractorClient) extractor.DependencyExtractorServer {
	return &extractorService{
		client: client,
	}
}

type extractorService struct {
	client extractor.DependencyExtractorClient
}

func (e *extractorService) Match(ctx context.Context, request *extractor.MatchRequest) (*extractor.MatchResponse, error) {
	return e.client.Match(ctx, request)
}

func (e *extractorService) Extract(ctx context.Context, request *extractor.ExtractRequest) (*extractor.ExtractResponse, error) {
	return e.client.Extract(ctx, request)
}

var _ extractor.DependencyExtractorServer = &extractorService{}
