package v1beta_test

import (
	"github.com/depscloud/depscloud/tracker/internal/graphstore/v1beta"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphData_TableName(t *testing.T) {
	{
		data := &v1beta.GraphData{
			CollectionName: "",
		}
		require.Equal(t, "graph_data", data.TableName())
	}

	{
		data := &v1beta.GraphData{
			CollectionName: "my_collection",
		}
		require.Equal(t, "my_collection", data.TableName())
	}
}
