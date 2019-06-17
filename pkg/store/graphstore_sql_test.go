package store_test

import (
	"database/sql"
	"testing"

	"github.com/deps-cloud/dts/pkg/store"
	_ "github.com/mattn/go-sqlite3"
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
		{GraphItemType: "node", K1: k1, K2: k1, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "node", K1: k2, K2: k2, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "node", K1: k3, K2: k3, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "node", K1: k4, K2: k4, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "node", K1: k5, K2: k5, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "node", K1: k6, K2: k6, Encoding: 0, GraphItemData: generateData()},

		{GraphItemType: "edge", K1: k1, K2: k2, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "edge", K1: k2, K2: k3, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "edge", K1: k2, K2: k4, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "edge", K1: k3, K2: k5, Encoding: 0, GraphItemData: generateData()},
		{GraphItemType: "edge", K1: k4, K2: k6, Encoding: 0, GraphItemData: generateData()},
	}

	rwdb, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.Nil(t, err)

	rodb, err := sql.Open("sqlite3", "file::memory:?cache=shared&mode=ro")
	require.Nil(t, err)

	graphStore, err := store.NewSQLGraphStore(rwdb, rodb)
	require.Nil(t, err)

	err = graphStore.Put(data)
	require.Nil(t, err)

	downstream, err := graphStore.FindDownstream(k2, []string{ "edge" })
	require.Nil(t, err)

	upstream, err := graphStore.FindUpstream(k2, []string{ "edge" })
	require.Nil(t, err)

	require.Len(t, downstream, 1)
	require.Len(t, upstream, 2)

	require.Equal(t, downstream[0].K1, k1)

	require.Equal(t, upstream[0].K1, k3)
	require.Equal(t, upstream[1].K1, k4)

	keys := make([]*store.PrimaryKey, 0, len(data))
	for _, data := range data {
		keys = append(keys, &store.PrimaryKey{
			GraphItemType: data.GraphItemType,
			K1:            data.K1,
			K2:            data.K2,
		})
	}

	err = graphStore.Delete(keys)
	require.Nil(t, err)
}
