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
	"github.com/depscloud/depscloud/internal/logger"
	"github.com/depscloud/depscloud/internal/mux"
	"github.com/depscloud/depscloud/internal/v"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/urfave/cli/v2"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
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
	version := v.Info{Version: version, Commit: commit, Date: date}

	loggerConfig, loggerFlags := logger.WithFlags(logger.DefaultConfig())
	serverConfig, serverFlags := mux.WithFlags(mux.DefaultConfig(version))

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

	flags := make([]cli.Flag, 0)
	flags = append(flags, loggerFlags...)
	flags = append(flags, serverFlags...)
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
			log := logger.MustGetLogger(loggerConfig)
			ctx := logger.ToContext(c.Context, log)

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
			moduleService := tracker.NewModuleServiceClient(trackerConn)
			dependencyService := tracker.NewDependencyServiceClient(trackerConn)
			extractorService := extractor.NewDependencyExtractorClient(extractorConn)
			searchService := tracker.NewSearchServiceClient(trackerConn)

			serverConfig.Checks = checks.Checks(extractorService, sourceService, moduleService)
			serverConfig.Endpoints = []mux.ServerEndpoint{
				func(ctx context.Context, grpcServer *grpc.Server, httpServer *http.ServeMux) {
					gatewayMux := runtime.NewServeMux()

					tracker.RegisterSourceServiceServer(grpcServer, proxies.NewSourceServiceProxy(sourceService))
					_ = tracker.RegisterSourceServiceHandlerClient(ctx, gatewayMux, sourceService)

					tracker.RegisterModuleServiceServer(grpcServer, proxies.NewModuleServiceProxy(moduleService))
					_ = tracker.RegisterModuleServiceHandlerClient(ctx, gatewayMux, moduleService)

					tracker.RegisterDependencyServiceServer(grpcServer, proxies.NewDependencyServiceProxy(dependencyService))
					_ = tracker.RegisterDependencyServiceHandlerClient(ctx, gatewayMux, dependencyService)

					extractor.RegisterDependencyExtractorServer(grpcServer, proxies.NewExtractorServiceProxy(extractorService))
					_ = extractor.RegisterDependencyExtractorHandlerClient(ctx, gatewayMux, extractorService)

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
				},
			}

			server := mux.NewServer(serverConfig)
			return server.Serve(ctx)
		},
	}

	_ = app.Run(os.Args)
}
