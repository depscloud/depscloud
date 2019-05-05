package store_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mjpitz/dts/pkg/store"
	"github.com/stretchr/testify/require"
)

var (
	k1 = []byte("1001")
	k2 = []byte("2002")
	k3 = []byte("3003")
	k4 = []byte("4004")
	k5 = []byte("5005")
	k6 = []byte("6006")
)

func generateData() []byte {
	return make([]byte, 0)
}

func TestNewSQLGraphStore_sqlite(t *testing.T) {
	data := []*store.GraphItem{
		{ GraphItemType: "node", K1: k1, K2: k1, Version: 1, Encoding: 0, GraphItemData: generateData() },
		{ GraphItemType: "node", K1: k2, K2: k2, Version: 1, Encoding: 0, GraphItemData: generateData() },
		{ GraphItemType: "node", K1: k3, K2: k3, Version: 1, Encoding: 0, GraphItemData: generateData() },
		{ GraphItemType: "node", K1: k4, K2: k4, Version: 1, Encoding: 0, GraphItemData: generateData() },
		{ GraphItemType: "node", K1: k5, K2: k5, Version: 1, Encoding: 0, GraphItemData: generateData() },
		{ GraphItemType: "node", K1: k6, K2: k6, Version: 1, Encoding: 0, GraphItemData: generateData() },

		{ GraphItemType: "edge", K1: k1, K2: k2, Version: 1, Encoding: 0, GraphItemData: generateData() },
		{ GraphItemType: "edge", K1: k2, K2: k3, Version: 1, Encoding: 0, GraphItemData: generateData() },
		{ GraphItemType: "edge", K1: k2, K2: k4, Version: 2, Encoding: 0, GraphItemData: generateData() },
		{ GraphItemType: "edge", K1: k3, K2: k5, Version: 1, Encoding: 0, GraphItemData: generateData() },
		{ GraphItemType: "edge", K1: k4, K2: k6, Version: 1, Encoding: 0, GraphItemData: generateData() },
	}

	db, err := sql.Open("sqlite3", ":memory:")
	require.Nil(t, err)

	graphStore, err := store.NewSQLGraphStore(db)
	require.Nil(t, err)

	for _, item := range data {
		err := graphStore.Put(item)
		require.Nil(t, err)
	}

	downstream, err := graphStore.FindDownstream(k2)
	require.Nil(t, err)

	upstream, err := graphStore.FindUpstream(k2)
	require.Nil(t, err)

	require.Len(t, downstream, 1)
	require.Len(t, upstream, 2)

	require.Equal(t, downstream[0].K1, k1)

	require.Equal(t, upstream[0].K1, k3)
	require.Equal(t, upstream[1].K1, k4)
}
