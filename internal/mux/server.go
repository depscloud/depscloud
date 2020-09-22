package mux

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"

	"github.com/mjpitz/go-gracefully/check"
	"github.com/mjpitz/go-gracefully/health"
	"github.com/mjpitz/go-gracefully/state"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rs/cors"

	"github.com/sirupsen/logrus"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Config struct {
	context.Context

	BindAddressHTTP string
	BindAddressGRPC string

	Checks []check.Check

	TLSConfig *TLSConfig
}

func registerHealth(grpcServer *grpc.Server, httpServer *http.ServeMux, config *Config) {
	monitor := health.NewMonitor(config.Checks...)
	reports, unsubscribe := monitor.Subscribe()
	stopCh := config.Context.Done()

	healthCheck := grpchealth.NewServer()

	go func() {
		defer unsubscribe()

		for {
			select {
			case <-stopCh:
				return
			case report := <-reports:
				if report.Check == nil {
					if report.Result.State == state.Outage {
						healthCheck.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
					} else {
						healthCheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
					}
				}
			}
		}
	}()

	httpServer.HandleFunc("/healthz", health.HandlerFunc(monitor))
	healthpb.RegisterHealthServer(grpcServer, healthCheck)
	_ = monitor.Start(config.Context)
}

func registerMetrics(httpServer *http.ServeMux) {
	httpServer.Handle("/metrics", promhttp.Handler())
}

func Serve(grpcServer *grpc.Server, httpServer http.Handler, config *Config) error {
	// TODO setup proper shutdown handlers

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.ProtoMajor == 2 &&
			strings.HasPrefix(request.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(writer, request)
		} else {
			httpServer.ServeHTTP(writer, request)
		}
	})

	registerHealth(grpcServer, httpMux, config)
	registerMetrics(httpMux)

	corsMux := cors.Default().Handler(httpMux)
	h2cMux := h2c.NewHandler(corsMux, &http2.Server{})

	var grpcListener net.Listener
	var grpcErr error

	var httpListener net.Listener
	var httpErr error

	tlsConfig, err := LoadTLSConfig(config.TLSConfig)
	if err != nil {
		return err
	}

	if tlsConfig != nil {
		httpListener, httpErr = tls.Listen("tcp", config.BindAddressHTTP, tlsConfig)
		grpcListener, grpcErr = tls.Listen("tcp", config.BindAddressGRPC, tlsConfig)
	} else {
		httpListener, httpErr = net.Listen("tcp", config.BindAddressHTTP)
		grpcListener, grpcErr = net.Listen("tcp", config.BindAddressGRPC)
	}

	if httpErr != nil {
		return httpErr
	}
	if grpcErr != nil {
		return grpcErr
	}

	logrus.Infof("[runtime] starting http on %s", config.BindAddressHTTP)
	go http.Serve(httpListener, h2cMux)

	logrus.Infof("[runtime] starting grpc on %s", config.BindAddressGRPC)
	return grpcServer.Serve(grpcListener)
}
