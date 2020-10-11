package client

import "github.com/depscloud/api/v1alpha/tracker"

type httpClient struct {
	dependencies tracker.DependencyServiceClient
	modules      tracker.ModuleServiceClient
	sources      tracker.SourceServiceClient
	search       tracker.SearchServiceClient
	troubleshoot *httpTroubleshootClient
}

func (c *httpClient) Dependencies() tracker.DependencyServiceClient {
	return c.dependencies
}

func (c *httpClient) Modules() tracker.ModuleServiceClient {
	return c.modules
}

func (c *httpClient) Sources() tracker.SourceServiceClient {
	return c.sources
}

func (c *httpClient) Search() tracker.SearchServiceClient {
	return c.search
}

func (c *httpClient) Troubleshoot() *httpTroubleshootClient {
	return c.troubleshoot
}

var _ Client = &httpClient{}
