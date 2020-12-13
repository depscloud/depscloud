package checks

import (
	"context"
	"time"

	apiv1alpha "github.com/depscloud/api/v1alpha/store"
	apiv1beta "github.com/depscloud/api/v1beta/graphstore"
	"github.com/depscloud/depscloud/tracker/internal/services/v1alpha"

	"github.com/mjpitz/go-gracefully/check"
	"github.com/mjpitz/go-gracefully/state"
)

// Checks returns an array of all health checks for the system.
func Checks(
	v1betaGraphStore apiv1beta.GraphStoreClient,
	graphStore apiv1alpha.GraphStoreClient,
) []check.Check {
	return []check.Check{
		&check.Periodic{
			Metadata: check.Metadata{
				Name:   "graphstore-v1beta-read",
				Weight: 10,
			},
			Interval: time.Second * 5,
			Timeout:  time.Second * 5,
			RunFunc: func(ctx context.Context) (state.State, error) {
				_, err := v1betaGraphStore.List(ctx, &apiv1beta.ListRequest{
					PageSize: 1,
					Kind:     v1alpha.SourceType,
				})
				if err != nil {
					return state.Outage, err
				}

				return state.OK, nil
			},
		},
		&check.Periodic{
			Metadata: check.Metadata{
				Name:   "graphstore-v1alpha-read",
				Weight: 10,
			},
			Interval: time.Second * 5,
			Timeout:  time.Second * 5,
			RunFunc: func(ctx context.Context) (state.State, error) {
				_, err := graphStore.List(ctx, &apiv1alpha.ListRequest{
					Count: 1,
					Type:  v1alpha.SourceType,
				})
				if err != nil {
					return state.Outage, err
				}

				return state.OK, nil
			},
		},
	}
}
