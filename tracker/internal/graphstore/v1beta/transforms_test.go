package v1beta_test

import (
	"testing"

	"github.com/depscloud/api/v1beta/graphstore"
	"github.com/depscloud/depscloud/tracker/internal/graphstore/v1beta"

	"github.com/golang/protobuf/ptypes/any"

	"github.com/stretchr/testify/require"
)

var testNode = &graphstore.Node{
	Key: []byte("key"),
	Body: &any.Any{
		TypeUrl: "kind",
		Value:   []byte("data"),
	},
}

var testGraphItemNode = &v1beta.GraphData{
	K1:       []byte("key"),
	K2:       []byte("key"),
	K3:       []byte{},
	Kind:     "kind",
	Encoding: v1beta.EncodingUnspecified,
	Data:     []byte("data"),
}

var testEdge = &graphstore.Edge{
	FromKey: []byte("fromKey"),
	ToKey:   []byte("toKey"),
	Key:     []byte("key"),
	Body: &any.Any{
		TypeUrl: "kind",
		Value:   []byte("data"),
	},
}
var testGraphItemEdge = &v1beta.GraphData{
	K1:       []byte("fromKey"),
	K2:       []byte("toKey"),
	K3:       []byte("key"),
	Kind:     "kind",
	Encoding: v1beta.EncodingUnspecified,
	Data:     []byte("data"),
}

func Test_ConvertNode(t *testing.T) {
	item := v1beta.ConvertNode(testNode, v1beta.EncodingUnspecified)

	require.Equal(t, string(testGraphItemNode.K1), string(item.K1))
	require.Equal(t, string(testGraphItemNode.K2), string(item.K2))
	require.Equal(t, string(testGraphItemNode.K3), string(item.K3))
	require.Equal(t, testGraphItemNode.Kind, item.Kind)
	require.Equal(t, testGraphItemNode.Encoding, item.Encoding)
	require.Equal(t, string(testGraphItemNode.Data), string(item.Data))
}

func Test_ConvertEdge(t *testing.T) {
	item := v1beta.ConvertEdge(testEdge, v1beta.EncodingUnspecified)

	require.Equal(t, string(testGraphItemEdge.K1), string(item.K1))
	require.Equal(t, string(testGraphItemEdge.K2), string(item.K2))
	require.Equal(t, string(testGraphItemEdge.K3), string(item.K3))
	require.Equal(t, testGraphItemEdge.Kind, item.Kind)
	require.Equal(t, testGraphItemEdge.Encoding, item.Encoding)
	require.Equal(t, string(testGraphItemEdge.Data), string(item.Data))
}

func Test_ConvertGraphItem_Node(t *testing.T) {
	node, edge, err := v1beta.ConvertGraphData(testGraphItemNode)
	require.Nil(t, err)
	require.Nil(t, edge)

	require.Equal(t, string(testNode.GetKey()), string(node.GetKey()))
	require.Equal(t, testNode.GetBody().GetTypeUrl(), node.GetBody().GetTypeUrl())
	require.Equal(t, string(testNode.GetBody().GetValue()), string(node.GetBody().GetValue()))
}

func Test_ConvertGraphItem_Edge(t *testing.T) {
	node, edge, err := v1beta.ConvertGraphData(testGraphItemEdge)
	require.Nil(t, err)
	require.Nil(t, node)

	require.Equal(t, string(testEdge.GetFromKey()), string(edge.GetFromKey()))
	require.Equal(t, string(testEdge.GetToKey()), string(edge.GetToKey()))
	require.Equal(t, string(testEdge.GetKey()), string(edge.GetKey()))
	require.Equal(t, testEdge.GetBody().GetTypeUrl(), edge.GetBody().GetTypeUrl())
	require.Equal(t, string(testEdge.GetBody().GetValue()), string(edge.GetBody().GetValue()))
}
