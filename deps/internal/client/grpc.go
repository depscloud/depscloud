package client

import (
	"crypto/tls"
	"net/url"
	"strings"

	"github.com/depscloud/api/v1alpha/tracker"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func grpcClient(baseURL string) Client {
	isSecure, hostPort := translateBaseURL(baseURL)

	var credentialOption grpc.DialOption
	if isSecure {
		// TODO: eventually add support for mutual TLS certs
		tlsConfig := &tls.Config{}
		creds := credentials.NewTLS(tlsConfig)
		credentialOption = grpc.WithTransportCredentials(creds)
	} else {
		credentialOption = grpc.WithInsecure()
	}

	conn, err := grpc.Dial(hostPort, credentialOption)
	if err != nil {
		panic(err)
	}

	return &client{
		dependencies: tracker.NewDependencyServiceClient(conn),
		modules:      tracker.NewModuleServiceClient(conn),
		sources:      tracker.NewSourceServiceClient(conn),
		search:       tracker.NewSearchServiceClient(conn),
	}
}
