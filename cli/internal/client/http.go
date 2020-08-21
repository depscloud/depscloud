package client

import "net/http"

func httpClient(baseURL string) Client {
	httpClient := http.DefaultClient

	return &client{
		dependencies: &dependencyService{httpClient, baseURL},
		modules:      &moduleClient{httpClient, baseURL},
		sources:      &sourceClient{httpClient, baseURL},
		search:       nil,
	}
}
