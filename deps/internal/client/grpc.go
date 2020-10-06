package client

import (
	"net/url"
	"strings"

	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/internal/client"
)

func translateBaseURL(baseURL string) (bool, string) {
	tls := false
	uri, _ := url.Parse(baseURL)

	if uri.Scheme == "https" {
		tls = true
	}

	host := uri.Host
	if !strings.Contains(host, ":") {
		if tls {
			host = host + ":443"
		} else {
			host = host + ":80"
		}
	}

	return tls, host
}

func grpcDefaultClient(baseURL string) Client {
	isSecure, hostPort := translateBaseURL(baseURL)

	conn, err := client.Connect(&client.Config{
		Address:       hostPort,
		ServiceConfig: client.DefaultServiceConfig,
		LoadBalancer:  client.DefaultLoadBalancer,
		TLS:           isSecure,
		TLSConfig:     &client.TLSConfig{},
	})
	if err != nil {
		panic(err)
	}

	return &httpClient{
		dependencies: tracker.NewDependencyServiceClient(conn),
		modules:      tracker.NewModuleServiceClient(conn),
		sources:      tracker.NewSourceServiceClient(conn),
		search:       tracker.NewSearchServiceClient(conn),
	}
}
