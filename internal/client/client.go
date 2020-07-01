package client

import (
	"os"

	"github.com/depscloud/api/v1alpha/tracker"
)

const (
	VariableProtocol = "DEPSCLOUD_PROTOCOL"
	VariableBaseURL  = "DEPSCLOUD_BASE_URL"

	DefaultProtocol = "grpc"
	DefaultBaseURL  = "https://api.deps.cloud"
)

var (
	protocol = or(os.Getenv(VariableProtocol), DefaultProtocol)
	baseURL  = or(os.Getenv(VariableBaseURL), DefaultBaseURL)
)

func or(read, def string) string {
	if read == "" {
		return def
	}
	return read
}

func DefaultClient() Client {
	if protocol == "grpc" {
		return grpcClient(baseURL)
	}

	return httpClient(baseURL)
}

type Client interface {
	Dependencies() tracker.DependencyServiceClient
	Modules() tracker.ModuleServiceClient
	Sources() tracker.SourceServiceClient
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
