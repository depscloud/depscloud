package v1alpha

import (
	"context"

	"github.com/cenkalti/backoff/v4"

	"github.com/depscloud/api/v1alpha/store"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
)

// Retryable wraps calls to the graphStore with an exponential backoff
func Retryable(delegate store.GraphStoreClient, maxAttempts int) store.GraphStoreClient {
	if maxAttempts <= 1 {
		return delegate
	}

	return &retryingClient{
		delegate:    delegate,
		maxAttempts: maxAttempts,
	}
}

type retryingClient struct {
	delegate    store.GraphStoreClient
	maxAttempts int
}

func (r *retryingClient) attempt(maxAttempts int, call func() error) error {
	backoffConfig := backoff.NewExponentialBackOff()

	attempt := 0
	var lastError error

	return backoff.Retry(func() error {
		if attempt++; attempt > maxAttempts {
			return backoff.Permanent(errors.Wrap(lastError, "max retries exceeded"))
		}

		if err := call(); err != nil {
			lastError = err
			return err
		}

		return nil
	}, backoffConfig)
}

func (r *retryingClient) Put(ctx context.Context, in *store.PutRequest, opts ...grpc.CallOption) (*store.PutResponse, error) {
	var resp *store.PutResponse
	var interr error

	err := r.attempt(r.maxAttempts, func() error {
		resp, interr = r.delegate.Put(ctx, in, opts...)
		return interr
	})

	return resp, err
}

func (r *retryingClient) Delete(ctx context.Context, in *store.DeleteRequest, opts ...grpc.CallOption) (*store.DeleteResponse, error) {
	var resp *store.DeleteResponse
	var interr error

	err := r.attempt(r.maxAttempts, func() error {
		resp, interr = r.delegate.Delete(ctx, in, opts...)
		return interr
	})

	return resp, err
}

func (r *retryingClient) List(ctx context.Context, in *store.ListRequest, opts ...grpc.CallOption) (*store.ListResponse, error) {
	var resp *store.ListResponse
	var interr error

	err := r.attempt(r.maxAttempts, func() error {
		resp, interr = r.delegate.List(ctx, in, opts...)
		return interr
	})

	return resp, err
}

func (r *retryingClient) FindUpstream(ctx context.Context, in *store.FindRequest, opts ...grpc.CallOption) (*store.FindResponse, error) {
	var resp *store.FindResponse
	var interr error

	err := r.attempt(r.maxAttempts, func() error {
		resp, interr = r.delegate.FindUpstream(ctx, in, opts...)
		return interr
	})

	return resp, err
}

func (r *retryingClient) FindDownstream(ctx context.Context, in *store.FindRequest, opts ...grpc.CallOption) (*store.FindResponse, error) {
	var resp *store.FindResponse
	var interr error

	err := r.attempt(r.maxAttempts, func() error {
		resp, interr = r.delegate.FindDownstream(ctx, in, opts...)
		return interr
	})

	return resp, err
}

var _ store.GraphStoreClient = &retryingClient{}
