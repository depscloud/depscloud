package mux

import (
	"net/http"
	"strings"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/rs/cors"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// ServerOption provides a post-endpoint hook for configuring aspects about the server.
type ServerOption func(s *Server)

// WithMetrics enables reporting on standard HTTP metrics.
func WithMetrics() ServerOption {
	return func(s *Server) {
		grpc_prometheus.Register(s.grpc)

		mdlw := middleware.New(middleware.Config{
			Recorder: metrics.NewRecorder(metrics.Config{}),
		})
		s.http = std.Handler("", mdlw, s.http)
	}
}

// WithDualServe configures the HTTP server to handle both HTTP and gRPC requests.
func WithDualServe() ServerOption {
	return func(s *Server) {
		httpMux := s.http
		switchMux := http.NewServeMux()
		switchMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			if request.ProtoMajor == 2 &&
				strings.HasPrefix(request.Header.Get("Content-Type"), "application/grpc") {
				s.grpc.ServeHTTP(writer, request)
			} else {
				httpMux.ServeHTTP(writer, request)
			}
		})
		s.http = switchMux
	}
}

// WithCORS configures the server to support CORS requests.
func WithCORS() ServerOption {
	return func(s *Server) {
		s.http = cors.Default().Handler(s.http)
	}
}

// WithH2C configures the server to respond to HTTP/2 plaintext requests.
func WithH2C() ServerOption {
	return func(s *Server) {
		s.http = h2c.NewHandler(s.http, &http2.Server{})
	}
}
