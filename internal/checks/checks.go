package checks

import (
	"context"
	"net/http"
	"time"

	"github.com/depscloud/api/v1alpha/extractor"
	"github.com/depscloud/api/v1alpha/tracker"

	"github.com/mjpitz/go-gracefully/check"
	"github.com/mjpitz/go-gracefully/health"
	"github.com/mjpitz/go-gracefully/state"

	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Checks(
	dependencyExtractor extractor.DependencyExtractorClient,
	sourceService tracker.SourceServiceClient,
	moduleService tracker.ModuleServiceClient,
) []check.Check {
	return []check.Check{
		&check.Periodic{
			Metadata: check.Metadata{
				Name:   "extraction",
				Weight: 10,
			},
			Interval: time.Second * 5,
			Timeout:  time.Second * 5,
			RunFunc: func(ctx context.Context) (state.State, error) {
				_, err := dependencyExtractor.Match(ctx, &extractor.MatchRequest{})
				if err != nil {
					return state.Outage, err
				}
				return state.OK, nil
			},
		},
		&check.Periodic{
			Metadata: check.Metadata{
				Name:   "sources",
				Weight: 10,
			},
			Interval: time.Second * 5,
			Timeout:  time.Second * 5,
			RunFunc: func(ctx context.Context) (state.State, error) {
				_, err := sourceService.List(ctx, &tracker.ListRequest{})
				if err != nil {
					return state.Outage, err
				}
				return state.OK, nil
			},
		},
		&check.Periodic{
			Metadata: check.Metadata{
				Name:   "modules",
				Weight: 10,
			},
			Interval: time.Second * 5,
			Timeout:  time.Second * 5,
			RunFunc: func(ctx context.Context) (state.State, error) {
				_, err := moduleService.List(ctx, &tracker.ListRequest{})
				if err != nil {
					return state.Outage, err
				}
				return state.OK, nil
			},
		},
	}
}

func RegisterHealthCheck(ctx context.Context, httpMux *http.ServeMux, grpcServer *grpc.Server, checks []check.Check) {
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

	httpMux.HandleFunc("/healthz", health.HandlerFunc(monitor))
	healthpb.RegisterHealthServer(grpcServer, healthCheck)
	_ = monitor.Start(ctx)
}
