package checks

import (
	"context"
	"time"

	"github.com/depscloud/api/v1beta"

	"github.com/mjpitz/go-gracefully/check"
	"github.com/mjpitz/go-gracefully/state"
)

func Checks(
	extractionService v1beta.ManifestExtractionServiceClient,
	sourceService v1beta.SourceServiceClient,
	moduleService v1beta.ModuleServiceClient,
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
				_, err := extractionService.Match(ctx, &v1beta.MatchRequest{})
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
				_, err := sourceService.List(ctx, &v1beta.ListRequest{})
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
				_, err := moduleService.List(ctx, &v1beta.ListRequest{})
				if err != nil {
					return state.Outage, err
				}
				return state.OK, nil
			},
		},
	}
}
