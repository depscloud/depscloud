package checks

import (
	"context"
	"time"

	// TODO: move over to v1beta post v0.3.0
	"github.com/depscloud/api/v1alpha/extractor"
	"github.com/depscloud/api/v1alpha/tracker"

	"github.com/mjpitz/go-gracefully/check"
	"github.com/mjpitz/go-gracefully/state"
)

func Checks(
	extractionService extractor.DependencyExtractorClient,
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
				_, err := extractionService.Match(ctx, &extractor.MatchRequest{})
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
