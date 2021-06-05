package v1beta

import (
	"context"
)

// Driver represents a generic interface for storing a graph.
type Driver interface {
	Put(ctx context.Context, data []*GraphData) error
	Delete(ctx context.Context, data []*GraphData) error
	List(ctx context.Context, kind string, offset, limit int) ([]*GraphData, bool, error)
	NeighborsTo(ctx context.Context, toKeys []string) ([]*GraphData, error)
	NeighborsFrom(ctx context.Context, fromKeys []string) ([]*GraphData, error)
}
