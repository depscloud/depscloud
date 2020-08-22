package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"

	"github.com/gogo/protobuf/jsonpb"

	"google.golang.org/grpc"
)

type moduleClient struct {
	client  *http.Client
	baseURL string
}

func (m *moduleClient) List(ctx context.Context, in *tracker.ListRequest, opts ...grpc.CallOption) (*tracker.ListModuleResponse, error) {
	uri := fmt.Sprintf("%s/v1alpha/modules?page=%d&count=%d",
		m.baseURL,
		in.Page,
		in.Count)

	r, err := m.client.Get(uri)
	if err != nil {
		return nil, err
	}

	resp := &tracker.ListModuleResponse{}
	if err := jsonpb.Unmarshal(r.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *moduleClient) ListSources(ctx context.Context, in *schema.Module, opts ...grpc.CallOption) (*tracker.ListSourcesResponse, error) {
	uri := fmt.Sprintf("%s/v1alpha/modules/source?language=%s&organization=%s&module=%s",
		m.baseURL,
		url.QueryEscape(in.Language),
		url.QueryEscape(in.Organization),
		url.QueryEscape(in.Module))

	r, err := m.client.Get(uri)
	if err != nil {
		return nil, err
	}

	resp := &tracker.ListSourcesResponse{}
	if err := jsonpb.Unmarshal(r.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *moduleClient) ListManaged(ctx context.Context, in *schema.Source, opts ...grpc.CallOption) (*tracker.ListManagedResponse, error) {
	uri := fmt.Sprintf("%s/v1alpha/modules/managed?url=%s",
		m.baseURL,
		url.QueryEscape(in.Url))

	r, err := m.client.Get(uri)
	if err != nil {
		return nil, err
	}

	resp := &tracker.ListManagedResponse{}
	if err := jsonpb.Unmarshal(r.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

var _ tracker.ModuleServiceClient = &moduleClient{}
