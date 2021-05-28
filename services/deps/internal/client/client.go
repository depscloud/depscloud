package client

import (
	"os"

	"github.com/depscloud/api/v1beta"
)

const (
	VariableBaseURL = "DEPSCLOUD_BASE_URL"

	VariableAddress       = "DEPSCLOUD_ADDRESS"
	VariableServiceConfig = "DEPSCLOUD_SERVICE_CONFIG"
	VariableTLS           = "DEPSCLOUD_TLS"
	VariableCAPath        = "DEPSCLOUD_CA_PATH"
	VariableCertPath      = "DEPSCLOUD_CERT_PATH"
	VariableKeyPath       = "DEPSCLOUD_KEY_PATH"

	DefaultBaseURL = "https://api.deps.cloud"
)

var (
	baseURL = or(os.Getenv(VariableBaseURL), DefaultBaseURL)
)

func or(read, def string) string {
	if read == "" {
		return def
	}
	return read
}

func DefaultClient() Client {
	return grpcDefaultClient(baseURL)
}

type Client interface {
	Modules() v1beta.ModuleServiceClient
	Sources() v1beta.SourceServiceClient
	Traversal() v1beta.TraversalServiceClient
}

type internalClient struct {
	moduleService    v1beta.ModuleServiceClient
	sourceService    v1beta.SourceServiceClient
	traversalService v1beta.TraversalServiceClient
}

func (c *internalClient) Modules() v1beta.ModuleServiceClient {
	return c.moduleService
}

func (c *internalClient) Sources() v1beta.SourceServiceClient {
	return c.sourceService
}

func (c *internalClient) Traversal() v1beta.TraversalServiceClient {
	return c.traversalService
}

var _ Client = &internalClient{}
