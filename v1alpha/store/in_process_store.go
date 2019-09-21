package store

import (
	"context"

	"google.golang.org/grpc"
)

// NewInProcessGraphStoreClient constructs a thin shim allowing a gRPC service to be called from in memory
func NewInProcessGraphStoreClient(server GraphStoreServer) GraphStoreClient {
	return &inProcessGraphStoreClient{
		server: server,
	}
}

type inProcessGraphStoreClient struct {
	server GraphStoreServer
}

var _ GraphStoreClient = &inProcessGraphStoreClient{}

func (c *inProcessGraphStoreClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error) {
	return c.server.Put(ctx, in)
}

func (c *inProcessGraphStoreClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	return c.server.Delete(ctx, in)
}

func (c *inProcessGraphStoreClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	return c.server.List(ctx, in)
}

func (c *inProcessGraphStoreClient) FindUpstream(ctx context.Context, in *FindRequest, opts ...grpc.CallOption) (*FindResponse, error) {
	return c.server.FindUpstream(ctx, in)
}

func (c *inProcessGraphStoreClient) FindDownstream(ctx context.Context, in *FindRequest, opts ...grpc.CallOption) (*FindResponse, error) {
	return c.server.FindDownstream(ctx, in)
}
