package client

import "net/http"

func httpDefaltClient(baseURL string) Client {
	client := http.DefaultClient

	return &httpClient{
		dependencies: &httpDependencyService{client, baseURL},
		modules:      &httpModuleClient{client, baseURL},
		sources:      &httpSourceClient{client, baseURL},
		troubleshoot: &httpTroubleshootClient{client, baseURL},
		search:       nil,
	}
}
