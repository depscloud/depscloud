package v1beta

import (
	"context"
	"database/sql"
	"encoding/base32"
	"encoding/binary"
	"time"

	"github.com/depscloud/api/v1beta/graphstore"
)

// GraphStoreServer encapsulates the logic for storing a graph using a generic driver.
type GraphStoreServer struct {
	Driver Driver
}

// Put ...
func (s *GraphStoreServer) Put(ctx context.Context, request *graphstore.PutRequest) (*graphstore.PutResponse, error) {
	length := len(request.GetNodes()) + len(request.GetEdges())
	items := make([]*GraphData, 0, length)

	// add nodes before adding edges
	// otherwise edges may point to nodes that do not exist

	for _, node := range request.GetNodes() {
		items = append(items, ConvertNode(node, EncodingProtocolBuffers))
	}

	for _, edge := range request.GetEdges() {
		items = append(items, ConvertEdge(edge, EncodingProtocolBuffers))
	}

	if err := s.Driver.Put(ctx, items); err != nil {
		return nil, err
	}

	return &graphstore.PutResponse{}, nil
}

// Delete ...
func (s *GraphStoreServer) Delete(ctx context.Context, request *graphstore.DeleteRequest) (*graphstore.DeleteResponse, error) {
	length := len(request.GetNodes()) + len(request.GetEdges())
	items := make([]*GraphData, 0, length)

	// delete edges before deleting nodes
	// otherwise edges may point to nodes that do not exist

	dateDeleted := &sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	for _, edge := range request.GetEdges() {
		item := ConvertEdge(edge, EncodingProtocolBuffers)
		item.DateDeleted = dateDeleted
		items = append(items, item)
	}

	for _, node := range request.GetNodes() {
		item := ConvertNode(node, EncodingProtocolBuffers)
		item.DateDeleted = dateDeleted
		items = append(items, item)
	}

	if err := s.Driver.Delete(ctx, items); err != nil {
		return nil, err
	}

	return &graphstore.DeleteResponse{}, nil
}

func paginate(lastToken string, pageSize int) (int, string) {
	offset := 0

	if lastToken != "" {
		data, err := base32.StdEncoding.DecodeString(lastToken)
		if err != nil {
			return offset, ""
		}
		offset = int(binary.BigEndian.Uint64(data))
	}

	nextOffset := offset + pageSize

	nextToken := make([]byte, 8)
	binary.BigEndian.PutUint64(nextToken, uint64(nextOffset))

	return offset, base32.StdEncoding.EncodeToString(nextToken)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// List ...
func (s *GraphStoreServer) List(ctx context.Context, request *graphstore.ListRequest) (*graphstore.ListResponse, error) {
	limit := int(request.GetPageSize())
	if limit == 0 {
		limit = 10
	}
	limit = min(100, limit)

	offset, nextPageToken := paginate(request.GetPageToken(), limit)

	results, hasNext, err := s.Driver.List(ctx, request.GetKind(), offset, limit)
	if err != nil {
		return nil, err
	}

	if !hasNext {
		nextPageToken = ""
	}

	nodes := make([]*graphstore.Node, 0)
	edges := make([]*graphstore.Edge, 0)

	for _, item := range results {
		if item == nil {
			break
		}

		node, edge, err := ConvertGraphData(item)
		if err != nil {
			continue
		}

		if node != nil {
			nodes = append(nodes, node)
		} else {
			edges = append(edges, edge)
		}
	}

	return &graphstore.ListResponse{
		NextPageToken: nextPageToken,
		Nodes:         nodes,
		Edges:         edges,
	}, nil
}

func processNeighbors(neighborData []*GraphData, keyFn func(edge *graphstore.Edge) string) []*graphstore.Neighbor {
	results := make([]*graphstore.Neighbor, 0)
	index := make(map[string]int)

	for _, neighbor := range neighborData {
		node, edge, err := ConvertGraphData(neighbor)
		if err != nil {
			continue
		}

		if node != nil {
			key := string(node.GetKey())
			if _, ok := index[key]; !ok {
				pos := len(results)
				results = append(results, &graphstore.Neighbor{
					Node:  node,
					Edges: make([]*graphstore.Edge, 0),
				})
				index[key] = pos
			}
		} else {
			pos := index[keyFn(edge)]
			results[pos].Edges = append(results[pos].Edges, edge)
		}
	}

	return results
}

// Neighbors ...
func (s *GraphStoreServer) Neighbors(ctx context.Context, request *graphstore.NeighborsRequest) (*graphstore.NeighborsResponse, error) {
	node := request.GetNode()
	toNode := request.GetTo()
	fromNode := request.GetFrom()

	if node != nil {
		toNode = node
		fromNode = node
	}

	neighbors := make([]*graphstore.Neighbor, 0)
	if toNode != nil {
		item := ConvertNode(toNode, EncodingProtocolBuffers)

		neighborData, err := s.Driver.NeighborsTo(ctx, []string{item.K1})
		if err != nil {
			return nil, err
		}

		toNeighbors := processNeighbors(neighborData, func(edge *graphstore.Edge) string {
			return string(edge.GetFromKey())
		})

		neighbors = append(neighbors, toNeighbors...)
	}

	if fromNode != nil {
		item := ConvertNode(fromNode, EncodingProtocolBuffers)

		neighborData, err := s.Driver.NeighborsFrom(ctx, []string{item.K1})
		if err != nil {
			return nil, err
		}

		fromNeighbors := processNeighbors(neighborData, func(edge *graphstore.Edge) string {
			return string(edge.GetToKey())
		})

		neighbors = append(neighbors, fromNeighbors...)
	}

	return &graphstore.NeighborsResponse{
		Neighbors: neighbors,
	}, nil
}

// Traverse ...
func (s *GraphStoreServer) Traverse(traverseServer graphstore.GraphStore_TraverseServer) error {
	ctx, cancel := context.WithCancel(traverseServer.Context())
	defer cancel()
	done := ctx.Done()

	incoming := make(chan *graphstore.NeighborsRequest, 1)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				req, err := traverseServer.Recv()
				if err != nil || req.Cancel {
					cancel()
					return
				}

				incoming <- req.Request
			}
		}
	}()

	for {
		select {
		case <-done:
			return nil

		case req := <-incoming:
			resp, err := s.Neighbors(ctx, req)
			if err != nil {
				return err
			}

			err = traverseServer.Send(&graphstore.TraverseResponse{
				Request:  req,
				Response: resp,
			})
			if err != nil {
				return err
			}
		}
	}
}

var _ graphstore.GraphStoreServer = &GraphStoreServer{}
