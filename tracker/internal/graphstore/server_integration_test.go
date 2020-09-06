// +build integration

package graphstore_test

import (
	"testing"

	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/depscloud/tracker/internal/graphstore"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/stretchr/testify/require"
)

func TestNewSQLGraphStore_postgres(t *testing.T) {
	data := []*store.GraphItem{
		{GraphItemType: "node", K1: k1, K2: k1, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "node", K1: k2, K2: k2, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "node", K1: k3, K2: k3, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "node", K1: k4, K2: k4, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "node", K1: k5, K2: k5, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "node", K1: k6, K2: k6, Encoding: 0, GraphItemData: generateData()},

		{GraphItemType: "edge", K1: k1, K2: k2, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "edge", K1: k2, K2: k3, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "edge", K1: k2, K2: k4, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "edge", K1: k4, K2: k6, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "edge", K1: k3, K2: k5, K3: k1, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "edge", K1: k3, K2: k5, K3: k2, Encoding: 0, GraphItemData: generateData()},
	}

	// TODO: Figure a way to pass the driver and address as parameters from Github actions
	rwdb, err := sqlx.Open("pgx", "postgres://user:password@depscloud-postgresql:5432/depscloud")
	require.Nil(t, err)

	rodb := rwdb

	statements, err := graphstore.DefaultStatementsFor("pgx")
	require.Nil(t, err)

	graphStore, err := graphstore.NewSQLGraphStore(rwdb, rodb, statements)
	require.Nil(t, err)

	_, err = graphStore.Put(nil, &store.PutRequest{
		Items: data,
	})
	require.Nil(t, err)

	response, err := graphStore.List(nil, &store.ListRequest{
		Page:  1,
		Count: 10,
		Type:  "edge",
	})
	require.Nil(t, err)
	require.Len(t, response.Items, 6)

	downstream, err := graphStore.FindDownstream(nil, &store.FindRequest{
		Keys:      [][]byte{k2},
		EdgeTypes: []string{"edge"},
		NodeTypes: []string{"node"},
	})
	require.Nil(t, err)

	upstream, err := graphStore.FindUpstream(nil, &store.FindRequest{
		Keys:      [][]byte{k2},
		EdgeTypes: []string{"edge"},
		NodeTypes: []string{"node"},
	})
	require.Nil(t, err)

	require.Len(t, downstream.Pairs, 1)
	require.Len(t, upstream.Pairs, 2)

	require.Equal(t, downstream.Pairs[0].Node.K1, k1)
	require.Equal(t, downstream.Pairs[0].Edge.K1, k1)
	require.Equal(t, downstream.Pairs[0].Edge.K2, k2)

	require.Equal(t, upstream.Pairs[0].Node.K1, k3)
	require.Equal(t, upstream.Pairs[0].Edge.K1, k2)
	require.Equal(t, upstream.Pairs[0].Edge.K2, k3)

	require.Equal(t, upstream.Pairs[1].Node.K1, k4)
	require.Equal(t, upstream.Pairs[1].Edge.K1, k2)
	require.Equal(t, upstream.Pairs[1].Edge.K2, k4)

	// Tests for multiple edges between nodes
	upstreamNodeK3, err := graphStore.FindUpstream(nil, &store.FindRequest{
		Keys:      [][]byte{k3},
		EdgeTypes: []string{"edge"},
		NodeTypes: []string{"node"},
	})
	require.Nil(t, err)

	require.Len(t, upstreamNodeK3.Pairs, 2)
	require.Equal(t, upstreamNodeK3.Pairs[0].GetEdge().GetK1(), k3)
	require.Equal(t, upstreamNodeK3.Pairs[0].GetEdge().GetK2(), k5)
	require.Equal(t, upstreamNodeK3.Pairs[0].GetEdge().GetK3(), k1)
	require.Equal(t, upstreamNodeK3.Pairs[1].GetEdge().GetK1(), k3)
	require.Equal(t, upstreamNodeK3.Pairs[1].GetEdge().GetK2(), k5)
	require.Equal(t, upstreamNodeK3.Pairs[1].GetEdge().GetK3(), k2)

	_, err = graphStore.Delete(nil, &store.DeleteRequest{
		Items: data,
	})
	require.Nil(t, err)
}
