package http

import (
	"net/http"
	"os"

	"github.com/depscloud/api/v1alpha/tracker"
)

func DefaultClient() Client {
	baseURL := os.Getenv("DEPSCLOUD_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deps.cloud"
	}
	return NewClient(baseURL)
}

type Client interface {
	Dependencies() tracker.DependencyServiceClient
	Modules() tracker.ModuleServiceClient
	Sources() tracker.SourceServiceClient
}

func NewClient(baseURL string) Client {
	httpClient := http.DefaultClient

	return &client{
		dependencies: &dependencyService{httpClient, baseURL},
		modules:      &moduleClient{httpClient, baseURL},
		sources:      &sourceClient{httpClient, baseURL},
	}
}

type client struct {
	dependencies tracker.DependencyServiceClient
	modules      tracker.ModuleServiceClient
	sources      tracker.SourceServiceClient
}

func (c *client) Dependencies() tracker.DependencyServiceClient {
	return c.dependencies
}

func (c *client) Modules() tracker.ModuleServiceClient {
	return c.modules
}

func (c *client) Sources() tracker.SourceServiceClient {
	return c.sources
}

var _ Client = &client{}
