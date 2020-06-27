package http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/depscloud/api/v1alpha/tracker"

	"github.com/gogo/protobuf/jsonpb"

	"google.golang.org/grpc"
)

type dependencyService struct {
	client  *http.Client
	baseURL string
}

func (d *dependencyService) ListDependents(ctx context.Context, in *tracker.DependencyRequest, opts ...grpc.CallOption) (*tracker.ListDependentsResponse, error) {
	uri := fmt.Sprintf("%s/v1alpha/graph/%s/dependents?organization=%s&module=%s",
		d.baseURL,
		url.QueryEscape(in.Language),
		url.QueryEscape(in.Organization),
		url.QueryEscape(in.Module))

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

func (d *dependencyService) ListDependencies(ctx context.Context, in *tracker.DependencyRequest, opts ...grpc.CallOption) (*tracker.ListDependenciesResponse, error) {
	uri := fmt.Sprintf("%s/v1alpha/graph/%s/dependencies?organization=%s&module=%s",
		d.baseURL,
		url.QueryEscape(in.Language),
		url.QueryEscape(in.Organization),
		url.QueryEscape(in.Module))

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

var _ tracker.DependencyServiceClient = &dependencyService{}
