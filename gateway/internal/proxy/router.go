package proxy

import (
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Backend defines a backend service to route to.
type Backend struct {
	ClientConn      *grpc.ClientConn
	RegisterService func(server *grpc.Server)
}

type upstream struct {
	clientConn *grpc.ClientConn
	services   map[string]grpc.ServiceInfo
}

// NewRouter constructs a Router using the provided backends. Routing decisions
// are based on first match.
func NewRouter(backends ...*Backend) (*Router, error) {
	upstreams := make([]*upstream, 0, len(backends))

	for _, backend := range backends {
		fakeServer := grpc.NewServer()
		backend.RegisterService(fakeServer)

		services := fakeServer.GetServiceInfo()
		upstreams = append(upstreams, &upstream{
			clientConn: backend.ClientConn,
			services:   services,
		})
	}

	return &Router{upstreams}, nil
}

// Router maintains mappings of services and their upstream channels.
type Router struct {
	upstreams []*upstream
}

// Route determines which upstream a given request should be sent to.
func (s *Router) Route(fullMethodName string) (*grpc.ClientConn, error) {
	parts := strings.Split(fullMethodName, "/")
	serviceName := parts[1]

	for _, upstream := range s.upstreams {
		if _, ok := upstream.services[serviceName]; ok {
			return upstream.clientConn, nil
		}
	}

	return nil, status.Error(codes.Unimplemented, "")
}
