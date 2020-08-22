package checks

import (
	"context"
	"time"

	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/depscloud/tracker/internal/types"

	"github.com/mjpitz/go-gracefully/check"
	"github.com/mjpitz/go-gracefully/health"
	"github.com/mjpitz/go-gracefully/state"

	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Checks returns an array of all health checks for the system.
func Checks(
	graphStore store.GraphStoreClient,
) []check.Check {
	return []check.Check{
		&check.Periodic{
			Metadata: check.Metadata{
				Name:   "graphstore",
				Weight: 10,
			},
			Interval: time.Second * 5,
			Timeout:  time.Second * 5,
			RunFunc: func(ctx context.Context) (state.State, error) {
				_, err := graphStore.List(ctx, &store.ListRequest{
					Count: 1,
					Type:  types.SourceType,
				})
				if err != nil {
					return state.Outage, err
				}

				return state.OK, nil
			},
		},
	}
}

// RegisterHealthCheck sets up the grpc server and associated monitor using the checks.
func RegisterHealthCheck(ctx context.Context, grpcServer *grpc.Server, checks []check.Check) {
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

	healthpb.RegisterHealthServer(grpcServer, healthCheck)
	_ = monitor.Start(ctx)
}
