package v1beta

import (
	"bytes"
	"encoding/base64"
	"time"

	"github.com/depscloud/api/v1beta/graphstore"

	"github.com/golang/protobuf/ptypes/any"
)

func encode(data []byte, encoding Encoding) []byte {
	switch encoding {
	case EncodingProtocolBuffers:
		return []byte(base64.StdEncoding.EncodeToString(data))
	default:
		return data
	}
}

func decode(data []byte, encoding Encoding) ([]byte, error) {
	switch encoding {
	case EncodingProtocolBuffers:
		return base64.StdEncoding.DecodeString(string(data))
	default:
		return data, nil
	}
}

// ConvertNode transforms a Node to a generalized graph item
func ConvertNode(node *graphstore.Node, encoding Encoding) *GraphData {
	return &GraphData{
		K1:           node.GetKey(),
		K2:           node.GetKey(),
		K3:           []byte{},
		Kind:         node.GetBody().GetTypeUrl(),
		Encoding:     encoding,
		Data:         encode(node.GetBody().GetValue(), encoding),
		DateDeleted:  nil,
		LastModified: time.Now(),
	}
}

// ConvertEdge transforms an Edge to a generalized graph item
func ConvertEdge(edge *graphstore.Edge, encoding Encoding) *GraphData {
	return &GraphData{
		K1:           edge.GetFromKey(),
		K2:           edge.GetToKey(),
		K3:           edge.GetKey(),
		Kind:         edge.GetBody().GetTypeUrl(),
		Encoding:     encoding,
		Data:         encode(edge.GetBody().GetValue(), encoding),
		DateDeleted:  nil,
		LastModified: time.Now(),
	}
}

// ConvertGraphData transforms itself back into a node or an edge
func ConvertGraphData(item *GraphData) (*graphstore.Node, *graphstore.Edge, error) {
	var node *graphstore.Node
	var edge *graphstore.Edge

	value, err := decode(item.Data, item.Encoding)
	if err != nil {
		return node, edge, err
	}

	body := &any.Any{
		TypeUrl: item.Kind,
		Value:   value,
	}

	if bytes.Equal(item.K1, item.K2) {
		node = &graphstore.Node{
			Key:  item.K1,
			Body: body,
		}
	} else {
		edge = &graphstore.Edge{
			FromKey: item.K1,
			ToKey:   item.K2,
			Key:     item.K3,
			Body:    body,
		}
	}

	return node, edge, nil
}
