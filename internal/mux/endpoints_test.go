package mux_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/depscloud/depscloud/internal/appconf"
	"github.com/depscloud/depscloud/internal/mux"

	"google.golang.org/grpc"
)

func testEndpoint(t *testing.T, endpoint mux.ServerEndpoint) {
	ctx := context.Background()
	grpcServer := grpc.NewServer()
	httpServer := http.NewServeMux()
	endpoint(ctx, grpcServer, httpServer)
}

func TestWithHealthEndpoint(t *testing.T) {
	testEndpoint(t, mux.WithHealthEndpoint())
}

func TestWithMetricsEndpoint(t *testing.T) {
	testEndpoint(t, mux.WithMetricsEndpoint())
}

func TestWithVersionEndpoint(t *testing.T) {
	testEndpoint(t, mux.WithVersionEndpoint(&appconf.V{}))
}

func TestWithReflectionEndpoint(t *testing.T) {
	testEndpoint(t, mux.WithReflectionEndpoint())
}
