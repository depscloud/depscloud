package mux

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/depscloud/depscloud/internal/logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

func newGRPC(log *zap.Logger) *grpc.Server {
	grpcOpts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			logger.StreamServerInterceptor(log),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			logger.UnaryServerInterceptor(log),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}

	grpc_prometheus.EnableHandlingTimeHistogram()

	return grpc.NewServer(grpcOpts...)
}

func newHTTP() *http.ServeMux {
	return http.NewServeMux()
}

// NewServer constructs a server given the provided configuration. By default,
// the server comes with the following:
//
//  * http server
//    * /health,/healthz endpoints
//    * /metrics endpoint
//    * /version endpoint
//    * prometheus metrics
//    * dual serve gRPC
//    * CORS
//    * H2C
//  * grpc server
//    * healthcheck service
//    * reflection service
//    * prometheus metrics
//
func NewServer(cfg *Config) *Server {
	endpoints := []ServerEndpoint{
		WithHealthEndpoint(cfg.Checks...),
		WithMetricsEndpoint(),
		WithVersionEndpoint(cfg.Version),
		WithReflectionEndpoint(),
	}
	endpoints = append(endpoints, cfg.Endpoints...)

	options := []ServerOption{
		WithMetrics(),
		WithDualServe(),
		WithCORS(),
		WithH2C(),
	}

	return &Server{
		mu:              &sync.Mutex{},
		http:            nil,
		httpBindAddress: fmt.Sprintf("0.0.0.0:%d", cfg.PortHTTP),
		grpc:            nil,
		grpcBindAddress: fmt.Sprintf("0.0.0.0:%d", cfg.PortGRPC),
		endpoints:       endpoints,
		options:         options,
		tlsConfig:       cfg.TLSConfig,
	}
}

// Server manages the various servers, endpoints, and options used to handle
// traffic. It drastically simplifies setting up and tearing down two servers.
type Server struct {
	mu          *sync.Mutex
	initialized bool

	http            http.Handler
	httpBindAddress string

	grpc            *grpc.Server
	grpcBindAddress string

	endpoints []ServerEndpoint
	options   []ServerOption

	tlsConfig *TLSConfig
}

// init is responsible for lazily creating the servers, configuring the
// endpoints, and setting up the various options. This block is protected
// by a mu to prevent concurrent initializations of the same server.
func (s *Server) init(root context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.initialized {
		return fmt.Errorf("already initialized")
	}

	log := logger.Extract(root)

	grpcServer := newGRPC(log)
	httpServer := newHTTP()

	for _, ep := range s.endpoints {
		ep(root, grpcServer, httpServer)
	}

	s.grpc = grpcServer
	s.http = httpServer

	for _, option := range s.options {
		option(s)
	}

	s.initialized = true
	return nil
}

// Serve initializes and boots the server. If the server has already been
// initialized, it returns an error.
func (s *Server) Serve(root context.Context) error {
	if err := s.init(root); err != nil {
		return err
	}

	log := logger.Extract(root)
	log.Info("initializing server")

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer signal.Stop(stop)

		<-stop
		log.Info("gracefully shutting down")
		go s.grpc.GracefulStop()

		<-stop
		log.Info("forcing termination")
		s.grpc.Stop()
	}()

	var httpListener net.Listener
	var httpErr error

	var grpcListener net.Listener
	var grpcErr error

	tlsConfig, err := LoadTLSConfig(s.tlsConfig)
	if err != nil {
		return err
	}

	if tlsConfig != nil {
		httpListener, httpErr = tls.Listen("tcp", s.httpBindAddress, tlsConfig)
		grpcListener, grpcErr = tls.Listen("tcp", s.grpcBindAddress, tlsConfig)
	} else {
		httpListener, httpErr = net.Listen("tcp", s.httpBindAddress)
		grpcListener, grpcErr = net.Listen("tcp", s.grpcBindAddress)
	}

	if httpErr != nil {
		return httpErr
	}
	if grpcErr != nil {
		return grpcErr
	}

	defer httpListener.Close()
	defer grpcListener.Close()

	log.Info("starting server",
		zap.String("protocol", "http"),
		zap.String("bind", s.httpBindAddress),
		zap.Bool("tls", tlsConfig != nil))

	go http.Serve(httpListener, s.http)

	log.Info("starting server",
		zap.String("protocol", "grpc"),
		zap.String("bind", s.grpcBindAddress),
		zap.Bool("tls", tlsConfig != nil))

	return s.grpc.Serve(grpcListener)
}
