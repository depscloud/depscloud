package v1beta_test

import (
	"context"
	"testing"

	"github.com/depscloud/api/v1beta/graphstore"
	"github.com/depscloud/depscloud/tracker/internal/graphstore/v1beta"

	"github.com/golang/protobuf/ptypes/any"

	"github.com/stretchr/testify/require"
)

func newNode(id string) *graphstore.Node {
	return &graphstore.Node{
		Key: []byte(id),
		Body: &any.Any{
			TypeUrl: "node",
			Value:   []byte{},
		},
	}
}

func newEdge(from, to, key string) *graphstore.Edge {
	return &graphstore.Edge{
		FromKey: []byte(from),
		ToKey:   []byte(to),
		Key:     []byte(key),
		Body: &any.Any{
			TypeUrl: "edge",
			Value:   []byte{},
		},
	}
}

func testServer(t *testing.T, storageDriver v1beta.Driver) {
	ctx := context.Background()

	store := &v1beta.GraphStoreServer{
		Driver: storageDriver,
	}

	// Put
	{
		_, err := store.Put(ctx, &graphstore.PutRequest{
			Nodes: []*graphstore.Node{
				newNode("a"),
				newNode("b"),
				newNode("c"),
				newNode("d"),
			},
			Edges: []*graphstore.Edge{
				newEdge("a", "b", ""),
				newEdge("a", "c", ""),
				newEdge("b", "d", ""),
				newEdge("c", "d", ""),
			},
		})
		require.Nil(t, err)
	}

	// List page 1
	nextPageToken := ""
	{
		resp, err := store.List(ctx, &graphstore.ListRequest{
			PageToken: nextPageToken,
			PageSize:  2,
			Kind:      "node",
		})
		require.Nil(t, err)

		nextPageToken = resp.GetNextPageToken()
		require.Equal(t, "AAAAAAAAAAAAE===", nextPageToken)
		require.Len(t, resp.GetEdges(), 0)
		require.Len(t, resp.GetNodes(), 2)
	}

	// List page 2
	{
		resp, err := store.List(ctx, &graphstore.ListRequest{
			PageToken: nextPageToken,
			PageSize:  2,
			Kind:      "node",
		})
		require.Nil(t, err)

		nextPageToken = resp.GetNextPageToken()
		require.Equal(t, "", nextPageToken)
		require.Len(t, resp.GetEdges(), 0)
		require.Len(t, resp.GetNodes(), 2)
	}

	// Neighbors from
	{
		resp, err := store.Neighbors(ctx, &graphstore.NeighborsRequest{
			From: newNode("a"),
		})
		require.Nil(t, err)

		neighbors := resp.GetNeighbors()
		require.Len(t, neighbors, 2)
		require.Equal(t, "b", string(neighbors[0].Node.Key))
		require.Equal(t, "c", string(neighbors[1].Node.Key))
	}

	// Neighbors to
	{
		resp, err := store.Neighbors(ctx, &graphstore.NeighborsRequest{
			To: newNode("d"),
		})
		require.Nil(t, err)

		neighbors := resp.GetNeighbors()
		require.Len(t, neighbors, 2)
		require.Equal(t, "b", string(neighbors[0].Node.Key))
		require.Equal(t, "c", string(neighbors[1].Node.Key))
	}

	// Neighbors
	{
		resp, err := store.Neighbors(ctx, &graphstore.NeighborsRequest{
			Node: newNode("b"),
		})
		require.Nil(t, err)

		neighbors := resp.GetNeighbors()
		require.Len(t, neighbors, 2)
		require.Equal(t, "a", string(neighbors[0].Node.Key))
		require.Equal(t, "d", string(neighbors[1].Node.Key))
	}

	// Delete
	{
		_, err := store.Delete(ctx, &graphstore.DeleteRequest{
			Nodes: []*graphstore.Node{
				newNode("a"),
				newNode("b"),
				newNode("c"),
				newNode("d"),
			},
			Edges: []*graphstore.Edge{
				newEdge("a", "b", ""),
				newEdge("a", "c", ""),
				newEdge("b", "d", ""),
				newEdge("c", "d", ""),
			},
		})
		require.Nil(t, err)
	}
}
