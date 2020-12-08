package mux_test

import (
	"context"
	"github.com/depscloud/depscloud/internal/mux"
	"github.com/depscloud/depscloud/internal/v"
	"google.golang.org/grpc"
	"net/http"
	"testing"
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
	testEndpoint(t, mux.WithVersionEndpoint(v.Info{}))
}

func TestWithReflectionEndpoint(t *testing.T) {
	testEndpoint(t, mux.WithReflectionEndpoint())
}
