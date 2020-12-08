package mux

import (
	"context"
	"net/http"

	"github.com/depscloud/depscloud/internal/v"

	"github.com/mjpitz/go-gracefully/check"
	"github.com/mjpitz/go-gracefully/health"
	"github.com/mjpitz/go-gracefully/state"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// ServerEndpoint provides applications with ways to configure the underlying grpc and http servers.
type ServerEndpoint func(ctx context.Context, grpcServer *grpc.Server, httpServer *http.ServeMux)

// WithHealthEndpoint sets up a server endpoint that responds to health requests.
func WithHealthEndpoint(checks ...check.Check) ServerEndpoint {
	return func(ctx context.Context, grpcServer *grpc.Server, httpServer *http.ServeMux) {
		monitor := health.NewMonitor(checks...)
		reports, unsubscribe := monitor.Subscribe()
		stopCh := ctx.Done()

		healthCheck := grpchealth.NewServer()

		go func() {
			defer unsubscribe()

			for {
				select {
				case <-stopCh:
					return
				case report := <-reports:
					if report.Check == nil {
						if report.Result.State == state.Outage {
							healthCheck.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
						} else {
							healthCheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
						}
					}
				}
			}
		}()

		handler := health.HandlerFunc(monitor)
		httpServer.HandleFunc("/healthz", handler)
		httpServer.HandleFunc("/health", handler)

		healthpb.RegisterHealthServer(grpcServer, healthCheck)
		_ = monitor.Start(ctx)
	}
}

// WithMetricsEndpoint sets up an endpoint that provides detailed metrics about the process.
func WithMetricsEndpoint() ServerEndpoint {
	return func(_ context.Context, _ *grpc.Server, httpServer *http.ServeMux) {
		httpServer.Handle("/metrics", promhttp.Handler())
	}
}

// WithVersionEndpoint sets up an endpoint that announces version metadata about the process.
func WithVersionEndpoint(version v.Info) ServerEndpoint {
	return func(_ context.Context, _ *grpc.Server, httpServer *http.ServeMux) {
		httpServer.HandleFunc("/version", func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(version.String()))
		})
	}
}

// WithReflectionEndpoint configures the servers with a discovery service.
func WithReflectionEndpoint() ServerEndpoint {
	return func(_ context.Context, grpcServer *grpc.Server, httpServer *http.ServeMux) {
		reflection.Register(grpcServer)
	}
}
