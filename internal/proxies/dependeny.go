package proxies

import (
	"context"

	"github.com/depscloud/api/v1alpha/tracker"
)

func NewDependencyServiceProxy(client tracker.DependencyServiceClient) tracker.DependencyServiceServer {
	return &dependencyService{
		client: client,
	}
}

type dependencyService struct {
	client tracker.DependencyServiceClient
}

func (d *dependencyService) ListDependents(ctx context.Context, request *tracker.DependencyRequest) (*tracker.ListDependentsResponse, error) {
	return d.client.ListDependents(ctx, request)
}

func (d *dependencyService) ListDependencies(ctx context.Context, request *tracker.DependencyRequest) (*tracker.ListDependenciesResponse, error) {
	return d.client.ListDependencies(ctx, request)
}

var _ tracker.DependencyServiceServer = &dependencyService{}
