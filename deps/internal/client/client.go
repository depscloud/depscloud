package client

import (
	"fmt"
	"os"
	"runtime"

	"github.com/depscloud/api/v1alpha/tracker"

	"github.com/sirupsen/logrus"
)

const (
	VariableProtocol = "DEPSCLOUD_PROTOCOL"
	VariableBaseURL  = "DEPSCLOUD_BASE_URL"

	VariableAddress       = "DEPSCLOUD_ADDRESS"
	VariableServiceConfig = "DEPSCLOUD_SERVICE_CONFIG"
	VariableTLS           = "DEPSCLOUD_TLS"
	VariableCAPath        = "DEPSCLOUD_CA_PATH"
	VariableCertPath      = "DEPSCLOUD_CERT_PATH"
	VariableKeyPath       = "DEPSCLOUD_KEY_PATH"

	DefaultProtocol = "grpc"
	DefaultBaseURL  = "https://api.deps.cloud"
)

var (
	protocol = or(os.Getenv(VariableProtocol), DefaultProtocol)
	baseURL  = or(os.Getenv(VariableBaseURL), DefaultBaseURL)
	systemOS = runtime.GOOS
	arch     = runtime.GOARCH
)

type SystemInfo struct {
	Protocol string
	BaseURL  string
	Os       string
	Arch     string
}

func (s SystemInfo) String() string {
	return fmt.Sprintf("{Protocol: %v, BaseURL: %v, OS: %v, Arch: %v}", s.Protocol, s.BaseURL, s.Os, s.Arch)
}

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

func GetSystemInfo() SystemInfo {
	return SystemInfo{Protocol: protocol, BaseURL: baseURL, Os: systemOS, Arch: arch}
}

type Client interface {
	Dependencies() tracker.DependencyServiceClient
	Modules() tracker.ModuleServiceClient
	Sources() tracker.SourceServiceClient
	Search() tracker.SearchServiceClient
}
