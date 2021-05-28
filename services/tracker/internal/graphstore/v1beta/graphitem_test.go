package v1beta_test

import (
	"testing"

	"github.com/depscloud/depscloud/services/tracker/internal/graphstore/v1beta"

	"github.com/stretchr/testify/require"
)

func TestGraphData_TableName(t *testing.T) {
	{
		data := &v1beta.GraphData{}
		require.Equal(t, "graph_data", data.TableName())
	}
}
