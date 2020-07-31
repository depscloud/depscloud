package graphstore

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"

	"github.com/depscloud/api/v1alpha/store"

	"google.golang.org/grpc"
)

const maxAttempts = 5

// Retryable wraps calls to the graphStore with an exponential backoff
func Retryable(delegate store.GraphStoreClient) store.GraphStoreClient {
	return &retryingClient{
		delegate: delegate,
	}
}

type retryingClient struct {
	delegate store.GraphStoreClient
}

func (r *retryingClient) attempt(maxAttempts int, call func() error) error {
	backoffConfig := backoff.NewExponentialBackOff()

	attempt := 0

	return backoff.Retry(func() error {
		if attempt++; attempt > maxAttempts {
			return backoff.Permanent(fmt.Errorf("max retries exceeded"))
		}

		if err := call(); err != nil {
			return err
		}

		return nil
	}, backoffConfig)
}

func (r *retryingClient) Put(ctx context.Context, in *store.PutRequest, opts ...grpc.CallOption) (*store.PutResponse, error) {
	var resp *store.PutResponse
	var interr error

	err := r.attempt(maxAttempts, func() error {
		resp, interr = r.delegate.Put(ctx, in, opts...)
		return interr
	})

	return resp, err
}

func (r *retryingClient) Delete(ctx context.Context, in *store.DeleteRequest, opts ...grpc.CallOption) (*store.DeleteResponse, error) {
	var resp *store.DeleteResponse
	var interr error

	err := r.attempt(maxAttempts, func() error {
		resp, interr = r.delegate.Delete(ctx, in, opts...)
		return interr
	})

	return resp, err
}

func (r *retryingClient) List(ctx context.Context, in *store.ListRequest, opts ...grpc.CallOption) (*store.ListResponse, error) {
	var resp *store.ListResponse
	var interr error

	err := r.attempt(maxAttempts, func() error {
		resp, interr = r.delegate.List(ctx, in, opts...)
		return interr
	})

	return resp, err
}

func (r *retryingClient) FindUpstream(ctx context.Context, in *store.FindRequest, opts ...grpc.CallOption) (*store.FindResponse, error) {
	var resp *store.FindResponse
	var interr error

	err := r.attempt(maxAttempts, func() error {
		resp, interr = r.delegate.FindUpstream(ctx, in, opts...)
		return interr
	})

	return resp, err
}

func (r *retryingClient) FindDownstream(ctx context.Context, in *store.FindRequest, opts ...grpc.CallOption) (*store.FindResponse, error) {
	var resp *store.FindResponse
	var interr error

	err := r.attempt(maxAttempts, func() error {
		resp, interr = r.delegate.FindDownstream(ctx, in, opts...)
		return interr
	})

	return resp, err
}

var _ store.GraphStoreClient = &retryingClient{}
