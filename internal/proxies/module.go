package proxies

import (
	"context"

	"github.com/deps-cloud/api/v1alpha/schema"
	"github.com/deps-cloud/api/v1alpha/tracker"
)

func NewModuleServiceProxy(client tracker.ModuleServiceClient) tracker.ModuleServiceServer {
	return &moduleService{
		client: client,
	}
}

type moduleService struct {
	client tracker.ModuleServiceClient
}

func (m *moduleService) List(ctx context.Context, request *tracker.ListRequest) (*tracker.ListModuleResponse, error) {
	return m.client.List(ctx, request)
}

func (m *moduleService) ListSources(ctx context.Context, module *schema.Module) (*tracker.ListSourcesResponse, error) {
	return m.client.ListSources(ctx, module)
}

func (m *moduleService) ListManaged(ctx context.Context, source *schema.Source) (*tracker.ListManagedResponse, error) {
	return m.client.ListManaged(ctx, source)
}

var _ tracker.ModuleServiceServer = &moduleService{}
