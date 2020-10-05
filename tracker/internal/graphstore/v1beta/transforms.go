package v1beta

import (
	"encoding/base64"
	"time"

	"github.com/depscloud/api/v1beta/graphstore"

	"github.com/golang/protobuf/ptypes/any"
)

func encode(data []byte, encoding Encoding) string {
	switch encoding {
	case EncodingProtocolBuffers:
		return base64.StdEncoding.EncodeToString(data)
	default:
		return string(data)
	}
}

func decode(data string, encoding Encoding) ([]byte, error) {
	switch encoding {
	case EncodingProtocolBuffers:
		return base64.StdEncoding.DecodeString(data)
	default:
		return []byte(data), nil
	}
}

// ConvertNode transforms a Node to a generalized graph item
func ConvertNode(node *graphstore.Node, encoding Encoding) *GraphData {
	return &GraphData{
		K1:           string(node.GetKey()),
		K2:           string(node.GetKey()),
		K3:           "",
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
		K1:           string(edge.GetFromKey()),
		K2:           string(edge.GetToKey()),
		K3:           string(edge.GetKey()),
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

	if item.K1 == item.K2 {
		node = &graphstore.Node{
			Key:  []byte(item.K1),
			Body: body,
		}
	} else {
		edge = &graphstore.Edge{
			FromKey: []byte(item.K1),
			ToKey:   []byte(item.K2),
			Key:     []byte(item.K3),
			Body:    body,
		}
	}

	return node, edge, nil
}
