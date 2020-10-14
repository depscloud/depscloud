package client

import (
	"context"
	"net/http"

	"github.com/depscloud/api/v1alpha/tracker"

	"github.com/gogo/protobuf/jsonpb"

	"google.golang.org/grpc"
)

type httpDependencyService struct {
	client  *http.Client
	baseURL string
}

func (d *httpDependencyService) ListDependents(ctx context.Context, in *tracker.DependencyRequest, opts ...grpc.CallOption) (*tracker.ListDependentsResponse, error) {
	uri := getUri("%s/v1alpha/graph/%s/dependents?organization=%s&module=%s",
		d.baseURL,
		in.Language,
		in.Organization,
		in.Module,
		in.Name)

	r, err := d.client.Get(uri)
	if err != nil {
		return nil, err
	}

	resp := &tracker.ListDependentsResponse{}
	if err := jsonpb.Unmarshal(r.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *httpDependencyService) ListDependencies(ctx context.Context, in *tracker.DependencyRequest, opts ...grpc.CallOption) (*tracker.ListDependenciesResponse, error) {
	uri := getUri("%s/v1alpha/graph/%s/dependencies?organization=%s&module=%s",
		d.baseURL,
		in.Language,
		in.Organization,
		in.Module,
		in.Name)

	r, err := d.client.Get(uri)
	if err != nil {
		return nil, err
	}

	resp := &tracker.ListDependenciesResponse{}
	if err := jsonpb.Unmarshal(r.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

var _ tracker.DependencyServiceClient = &httpDependencyService{}
