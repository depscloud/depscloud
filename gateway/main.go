package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/depscloud/api/swagger"
	"github.com/depscloud/api/v1alpha/extractor"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/gateway/internal/checks"
	"github.com/depscloud/depscloud/gateway/internal/proxies"
	"github.com/depscloud/depscloud/internal/client"
	"github.com/depscloud/depscloud/internal/mux"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"

	"golang.org/x/net/context"

	_ "google.golang.org/grpc/health"
)

// variables set during build using -X ldflag
var version string
var commit string
var date string

type gatewayConfig struct {
	httpPort int
	grpcPort int
}

func main() {
	version := mux.Version{Version: version, Commit: commit, Date: date}
	cfg := &gatewayConfig{
		httpPort: 8080,
		grpcPort: 8090,
	}

	tlsConfig := &mux.TLSConfig{}

	extractorConfig, extractorFlags := client.WithFlags("extractor", &client.Config{
		Address:       "extractor:8090",
		ServiceConfig: client.DefaultServiceConfig,
		LoadBalancer:  client.DefaultLoadBalancer,
		TLS:           false,
		TLSConfig:     &client.TLSConfig{},
	})

	trackerConfig, trackerFlags := client.WithFlags("tracker", &client.Config{
		Address:       "tracker:8090",
		ServiceConfig: client.DefaultServiceConfig,
		LoadBalancer:  client.DefaultLoadBalancer,
		TLS:           false,
		TLSConfig:     &client.TLSConfig{},
	})

	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "http-port",
			Aliases:     []string{"port"},
			Usage:       "the port to run http on",
			Value:       cfg.httpPort,
			Destination: &cfg.httpPort,
			EnvVars:     []string{"HTTP_PORT"},
		},
		&cli.IntFlag{
			Name:        "grpc-port",
			Usage:       "the port to run grpc on",
			Value:       cfg.grpcPort,
			Destination: &cfg.grpcPort,
			EnvVars:     []string{"GRPC_PORT"},
		},
		&cli.StringFlag{
			Name:        "tls-key",
			Usage:       "path to the file containing the TLS private key",
			Value:       tlsConfig.KeyPath,
			Destination: &tlsConfig.KeyPath,
			EnvVars:     []string{"TLS_KEY_PATH"},
		},
		&cli.StringFlag{
			Name:        "tls-cert",
			Usage:       "path to the file containing the TLS certificate",
			Value:       tlsConfig.CertPath,
			Destination: &tlsConfig.CertPath,
			EnvVars:     []string{"TLS_CERT_PATH"},
		},
		&cli.StringFlag{
			Name:        "tls-ca",
			Usage:       "path to the file containing the TLS certificate authority",
			Value:       tlsConfig.CAPath,
			Destination: &tlsConfig.CAPath,
			EnvVars:     []string{"TLS_CA_PATH"},
		},
	}

	flags = append(flags, extractorFlags...)
	flags = append(flags, trackerFlags...)

	app := &cli.App{
		Name:  "gateway",
		Usage: "an HTTP/gRPC proxy to backend services",
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Output version information",
				Action: func(c *cli.Context) error {
					fmt.Println(fmt.Sprintf("%s %s", c.Command.Name, version))
					return nil
				},
			},
		},
		Flags: flags,
		Action: func(c *cli.Context) error {
			grpcServer, httpServer := mux.DefaultServers()
			gatewayMux := runtime.NewServeMux()

			ctx := context.Background()

			extractorConn, err := client.Connect(extractorConfig)
			if err != nil {
				return err
			}
			defer extractorConn.Close()

			trackerConn, err := client.Connect(trackerConfig)
			if err != nil {
				return err
			}
			defer trackerConn.Close()

			sourceService := tracker.NewSourceServiceClient(trackerConn)
			tracker.RegisterSourceServiceServer(grpcServer, proxies.NewSourceServiceProxy(sourceService))
			_ = tracker.RegisterSourceServiceHandlerClient(ctx, gatewayMux, sourceService)

			moduleService := tracker.NewModuleServiceClient(trackerConn)
			tracker.RegisterModuleServiceServer(grpcServer, proxies.NewModuleServiceProxy(moduleService))
			_ = tracker.RegisterModuleServiceHandlerClient(ctx, gatewayMux, moduleService)

			dependencyService := tracker.NewDependencyServiceClient(trackerConn)
			tracker.RegisterDependencyServiceServer(grpcServer, proxies.NewDependencyServiceProxy(dependencyService))
			_ = tracker.RegisterDependencyServiceHandlerClient(ctx, gatewayMux, dependencyService)

			extractorService := extractor.NewDependencyExtractorClient(extractorConn)
			extractor.RegisterDependencyExtractorServer(grpcServer, proxies.NewExtractorServiceProxy(extractorService))
			_ = extractor.RegisterDependencyExtractorHandlerClient(ctx, gatewayMux, extractorService)

			searchService := tracker.NewSearchServiceClient(trackerConn)
			tracker.RegisterSearchServiceServer(grpcServer, proxies.NewSearchServiceProxy(searchService))

			httpServer.HandleFunc("/swagger/", func(writer http.ResponseWriter, request *http.Request) {
				assetPath := strings.TrimPrefix(request.URL.Path, "/swagger/")

				if len(assetPath) == 0 {
					if err := json.NewEncoder(writer).Encode(swagger.AssetNames()); err != nil {
						writer.WriteHeader(500)
					} else {
						writer.WriteHeader(200)
					}
					return
				}

				asset, err := swagger.Asset(assetPath)
				if err != nil {
					writer.WriteHeader(404)
					return
				}

				writer.WriteHeader(200)
				writer.Header().Set("Content-Type", "application/json")
				_, _ = writer.Write(asset)
			})

			httpServer.Handle("/", gatewayMux)

			return mux.Serve(grpcServer, httpServer, &mux.Config{
				Context:         c.Context,
				BindAddressHTTP: fmt.Sprintf("0.0.0.0:%d", cfg.httpPort),
				BindAddressGRPC: fmt.Sprintf("0.0.0.0:%d", cfg.grpcPort),
				Checks:          checks.Checks(extractorService, sourceService, moduleService),
				Version:         &version,
				TLSConfig:       tlsConfig,
			})
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
