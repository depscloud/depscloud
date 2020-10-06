package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/depscloud/api/v1alpha/tracker"

	"github.com/gogo/protobuf/jsonpb"

	"google.golang.org/grpc"
)

type httpSourceClient struct {
	client  *http.Client
	baseURL string
}

func (s *httpSourceClient) List(ctx context.Context, in *tracker.ListRequest, opts ...grpc.CallOption) (*tracker.ListSourceResponse, error) {
	uri := fmt.Sprintf("%s/v1alpha/sources?page=%d&count=%d",
		s.baseURL,
		in.Page,
		in.Count)

	r, err := s.client.Get(uri)
	if err != nil {
		return nil, err
	}

	resp := &tracker.ListSourceResponse{}
	if err := jsonpb.Unmarshal(r.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *httpSourceClient) Track(ctx context.Context, in *tracker.SourceRequest, opts ...grpc.CallOption) (*tracker.TrackResponse, error) {
	uri := fmt.Sprintf("%s/v1alpha/sources/track",
		s.baseURL)

	body, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	r, err := s.client.Post(uri, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp := &tracker.TrackResponse{}
	if err := jsonpb.Unmarshal(r.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

var _ tracker.SourceServiceClient = &httpSourceClient{}
