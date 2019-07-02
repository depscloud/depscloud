package service

import (
	"context"

	"github.com/deps-cloud/tracker/api"
	"github.com/deps-cloud/tracker/api/v1alpha/store"
)

// TraversalUtil encapsulates logic for traversing the dependency graph in the
// proper direction
type TraversalUtil struct {
	Graph     store.GraphStoreClient
	Direction api.Direction
}

// GetAdjacent retrieves adjacent nodes following the provided edges
func (tu *TraversalUtil) GetAdjacent(key []byte, edgeTypes []string) ([]*store.GraphItem, error) {
	var resp *store.FindResponse
	var err error

	if tu.Direction == api.Direction_UPSTREAM {
		resp, err = tu.Graph.FindUpstream(context.Background(), &store.FindRequest{
			Key:       key,
			EdgeTypes: edgeTypes,
		})
	} else {
		resp, err = tu.Graph.FindDownstream(context.Background(), &store.FindRequest{
			Key:       key,
			EdgeTypes: edgeTypes,
		})
	}

	if err != nil {
		return nil, err
	}

	result := make([]*store.GraphItem, 0, len(resp.GetPairs()))
	for _, pair := range resp.GetPairs() {
		result = append(result, pair.GetNode())
	}

	return result, nil
}
