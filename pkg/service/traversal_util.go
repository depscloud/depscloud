package service

import (
	"github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/pkg/store"
)

// TraversalUtil encapsulates logic for traversing the dependency graph in the
// proper direction
type TraversalUtil struct {
	Graph 	  store.GraphStore
	Direction api.Direction
}

// GetAdjacent retrieves adjacent nodes following the provided edges
func (tu *TraversalUtil) GetAdjacent(key []byte, edgeTypes []string) ([]*store.GraphItem, error) {
	if tu.Direction == api.Direction_UPSTREAM {
		return tu.Graph.FindUpstream(key, edgeTypes)
	}

	return tu.Graph.FindDownstream(key, edgeTypes)
}
