package graphstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/depscloud/api"
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/depscloud/tracker/internal/graphstore"

	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
)

type graphStoreClientFuncs struct {
	PutFunc            func(ctx context.Context, in *store.PutRequest, opts ...grpc.CallOption) (*store.PutResponse, error)
	DeleteFunc         func(ctx context.Context, in *store.DeleteRequest, opts ...grpc.CallOption) (*store.DeleteResponse, error)
	ListFunc           func(ctx context.Context, in *store.ListRequest, opts ...grpc.CallOption) (*store.ListResponse, error)
	FindUpstreamFunc   func(ctx context.Context, in *store.FindRequest, opts ...grpc.CallOption) (*store.FindResponse, error)
	FindDownstreamFunc func(ctx context.Context, in *store.FindRequest, opts ...grpc.CallOption) (*store.FindResponse, error)
}

func (c *graphStoreClientFuncs) Put(ctx context.Context, in *store.PutRequest, opts ...grpc.CallOption) (*store.PutResponse, error) {
	if c.PutFunc == nil {
		return nil, api.ErrUnimplemented
	}
	return c.PutFunc(ctx, in, opts...)
}

func (c *graphStoreClientFuncs) Delete(ctx context.Context, in *store.DeleteRequest, opts ...grpc.CallOption) (*store.DeleteResponse, error) {
	if c.DeleteFunc == nil {
		return nil, api.ErrUnimplemented
	}
	return c.DeleteFunc(ctx, in, opts...)
}

func (c *graphStoreClientFuncs) List(ctx context.Context, in *store.ListRequest, opts ...grpc.CallOption) (*store.ListResponse, error) {
	if c.ListFunc == nil {
		return nil, api.ErrUnimplemented
	}
	return c.ListFunc(ctx, in, opts...)
}

func (c *graphStoreClientFuncs) FindUpstream(ctx context.Context, in *store.FindRequest, opts ...grpc.CallOption) (*store.FindResponse, error) {
	if c.FindUpstreamFunc == nil {
		return nil, api.ErrUnimplemented
	}
	return c.FindUpstreamFunc(ctx, in, opts...)
}

func (c *graphStoreClientFuncs) FindDownstream(ctx context.Context, in *store.FindRequest, opts ...grpc.CallOption) (*store.FindResponse, error) {
	if c.FindDownstreamFunc == nil {
		return nil, api.ErrUnimplemented
	}
	return c.FindDownstreamFunc(ctx, in, opts...)
}

var _ store.GraphStoreClient = &graphStoreClientFuncs{}

const testAttempts = 2

func Test_retryablePut(t *testing.T) {
	counter := 0

	client := graphstore.Retryable(&graphStoreClientFuncs{
		PutFunc: func(ctx context.Context, in *store.PutRequest, opts ...grpc.CallOption) (*store.PutResponse, error) {
			counter++
			return nil, fmt.Errorf("Put failed")
		},
	}, testAttempts)

	resp, err := client.Put(context.Background(), &store.PutRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, "max retries exceeded: Put failed", err.Error())
	require.Equal(t, testAttempts, counter)
}

func Test_retryableDelete(t *testing.T) {
	counter := 0

	client := graphstore.Retryable(&graphStoreClientFuncs{
		DeleteFunc: func(ctx context.Context, in *store.DeleteRequest, opts ...grpc.CallOption) (*store.DeleteResponse, error) {
			counter++
			return nil, fmt.Errorf("Delete failed")
		},
	}, testAttempts)

	resp, err := client.Delete(context.Background(), &store.DeleteRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, "max retries exceeded: Delete failed", err.Error())
	require.Equal(t, testAttempts, counter)
}

func Test_retryableList(t *testing.T) {
	counter := 0

	client := graphstore.Retryable(&graphStoreClientFuncs{
		ListFunc: func(ctx context.Context, in *store.ListRequest, opts ...grpc.CallOption) (*store.ListResponse, error) {
			counter++
			return nil, fmt.Errorf("List failed")
		},
	}, testAttempts)

	resp, err := client.List(context.Background(), &store.ListRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, "max retries exceeded: List failed", err.Error())
	require.Equal(t, testAttempts, counter)
}

func Test_retryableFindUpstream(t *testing.T) {
	counter := 0

	client := graphstore.Retryable(&graphStoreClientFuncs{
		FindUpstreamFunc: func(ctx context.Context, in *store.FindRequest, opts ...grpc.CallOption) (*store.FindResponse, error) {
			counter++
			return nil, fmt.Errorf("FindUpstream failed")
		},
	}, testAttempts)

	resp, err := client.FindUpstream(context.Background(), &store.FindRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, "max retries exceeded: FindUpstream failed", err.Error())
	require.Equal(t, testAttempts, counter)
}

func Test_retryableFindDownstream(t *testing.T) {
	counter := 0

	client := graphstore.Retryable(&graphStoreClientFuncs{
		FindDownstreamFunc: func(ctx context.Context, in *store.FindRequest, opts ...grpc.CallOption) (*store.FindResponse, error) {
			counter++
			return nil, fmt.Errorf("FindDownstream failed")
		},
	}, testAttempts)

	resp, err := client.FindDownstream(context.Background(), &store.FindRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, "max retries exceeded: FindDownstream failed", err.Error())
	require.Equal(t, testAttempts, counter)
}
