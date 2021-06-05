package mux

import (
	"testing"

	"github.com/depscloud/depscloud/internal/appconf"

	"go.uber.org/zap"
)

func testOption(t *testing.T, option ServerOption) {
	server := NewServer(DefaultConfig(&appconf.V{}))
	server.grpc = newGRPC(zap.NewNop())
	server.http = newHTTP()
	option(server)
}

func TestWithCORS(t *testing.T) {
	testOption(t, WithCORS())
}

func TestWithDualServe(t *testing.T) {
	testOption(t, WithDualServe())
}

func TestWithH2C(t *testing.T) {
	testOption(t, WithH2C())
}

func TestWithMetrics(t *testing.T) {
	testOption(t, WithMetrics())
}
