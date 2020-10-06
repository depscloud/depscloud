package client

import (
	"os"

	"github.com/depscloud/api/v1alpha/tracker"

	"github.com/sirupsen/logrus"
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
		return grpcDefaultClient(baseURL)
	}

	logrus.Warnf("the HTTP api is deprecated, please migrate to gRPC")
	return httpDefaltClient(baseURL)
}

type Client interface {
	Dependencies() tracker.DependencyServiceClient
	Modules() tracker.ModuleServiceClient
	Sources() tracker.SourceServiceClient
	Search() tracker.SearchServiceClient
}
