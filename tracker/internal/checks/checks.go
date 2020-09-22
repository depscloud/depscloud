package checks

import (
	"context"
	"time"

	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/depscloud/tracker/internal/types"

	"github.com/mjpitz/go-gracefully/check"
	"github.com/mjpitz/go-gracefully/state"
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
