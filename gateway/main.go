package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/depscloud/api/swagger"
	"github.com/depscloud/api/v1alpha/extractor"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/depscloud/gateway/internal/checks"
	"github.com/depscloud/depscloud/gateway/internal/proxy"
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

func main() {
	version := v.Info{Version: version, Commit: commit, Date: date}

	loggerConfig, loggerFlags := logger.WithFlags(logger.DefaultConfig())
	serverConfig, serverFlags := mux.WithFlags(mux.DefaultConfig(version))

	defaultDialOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.ForceCodec(proxy.Codec()),
		),
	}

	extractorConfig, extractorFlags := client.WithFlags("extractor", &client.Config{
		Address:       "extractor:8090",
		ServiceConfig: client.DefaultServiceConfig,
		LoadBalancer:  client.DefaultLoadBalancer,
		TLS:           false,
		TLSConfig:     &client.TLSConfig{},
		DialOptions:   defaultDialOptions,
	})

	trackerConfig, trackerFlags := client.WithFlags("tracker", &client.Config{
		Address:       "tracker:8090",
		ServiceConfig: client.DefaultServiceConfig,
		LoadBalancer:  client.DefaultLoadBalancer,
		TLS:           false,
		TLSConfig:     &client.TLSConfig{},
		DialOptions:   defaultDialOptions,
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

			gatewayConn, err := client.Connect(&client.Config{
				Address: fmt.Sprintf("localhost:%d", serverConfig.PortGRPC),
				TLSConfig: &client.TLSConfig{
					CertPath: serverConfig.TLSConfig.CertPath,
					KeyPath:  serverConfig.TLSConfig.KeyPath,
					CAPath:   serverConfig.TLSConfig.CAPath,
				},
				DialOptions: defaultDialOptions,
			})
			if err != nil {
				return err
			}
			defer gatewayConn.Close()

			// Setup a router with various backends. To do this, we create fake gRPC servers that capture
			// upstream services for a given client connection. Service information is extracted from the
			// fake server. When a call comes in for a given service, it's routed to the first clientConn
			// with the registered service.
			//
			// Eventually, we could look at leveraging the reflection service to make this proxy fully
			// dynamic. Instead of having the client register relevant services, we can look them up using
			// the reflection service and setup a dynamic routing table.
			router, err := proxy.NewRouter([]*proxy.Backend{
				{
					ClientConn: extractorConn,
					RegisterService: func(server *grpc.Server) {
						extractor.RegisterDependencyExtractorServer(server, &extractor.UnimplementedDependencyExtractorServer{})
						v1beta.RegisterManifestExtractionServiceServer(server, &v1beta.UnimplementedManifestExtractionServiceServer{})
					},
				},
				{
					ClientConn: trackerConn,
					RegisterService: func(server *grpc.Server) {
						tracker.RegisterSourceServiceServer(server, &tracker.UnimplementedSourceServiceServer{})
						tracker.RegisterModuleServiceServer(server, &tracker.UnimplementedModuleServiceServer{})
						tracker.RegisterDependencyServiceServer(server, &tracker.UnimplementedDependencyServiceServer{})
						tracker.RegisterSearchServiceServer(server, &tracker.UnimplementedSearchServiceServer{})
						v1beta.RegisterManifestStorageServiceServer(server, &v1beta.UnimplementedManifestStorageServiceServer{})
						v1beta.RegisterModuleServiceServer(server, &v1beta.UnimplementedModuleServiceServer{})
						v1beta.RegisterSourceServiceServer(server, &v1beta.UnimplementedSourceServiceServer{})
						v1beta.RegisterTraversalServiceServer(server, &v1beta.UnimplementedTraversalServiceServer{})
					},
				},
			}...)
			if err != nil {
				return err
			}

			extractorService := extractor.NewDependencyExtractorClient(gatewayConn)
			//extractionService := v1beta.NewManifestExtractionServiceClient(gatewayConn)

			sourceService := tracker.NewSourceServiceClient(gatewayConn)
			moduleService := tracker.NewModuleServiceClient(gatewayConn)
			//dependencyService := tracker.NewDependencyServiceClient(gatewayConn)
			//storageService := v1beta.NewManifestStorageServiceClient(gatewayConn)
			//storageService := v1beta.NewManifestStorageServiceClient(gatewayConn)

			serverConfig.GRPC.ServerOptions = []grpc.ServerOption{
				grpc.CustomCodec(proxy.ServerCodec()),
				grpc.UnknownServiceHandler(proxy.UnknownServiceHandler(router)),
			}
			serverConfig.Checks = checks.Checks(extractorService, sourceService, moduleService)
			serverConfig.Endpoints = []mux.ServerEndpoint{
				func(ctx context.Context, grpcServer *grpc.Server, httpServer *http.ServeMux) {
					gatewayMux := runtime.NewServeMux()

					_ = extractor.RegisterDependencyExtractorHandler(ctx, gatewayMux, gatewayConn)
					_ = v1beta.RegisterManifestExtractionServiceHandler(ctx, gatewayMux, gatewayConn)

					_ = tracker.RegisterSourceServiceHandler(ctx, gatewayMux, gatewayConn)
					_ = tracker.RegisterModuleServiceHandler(ctx, gatewayMux, gatewayConn)
					_ = tracker.RegisterDependencyServiceHandler(ctx, gatewayMux, gatewayConn)

					_ = v1beta.RegisterManifestStorageServiceHandler(ctx, gatewayMux, gatewayConn)
					_ = v1beta.RegisterSourceServiceHandler(ctx, gatewayMux, gatewayConn)
					_ = v1beta.RegisterModuleServiceHandler(ctx, gatewayMux, gatewayConn)
					_ = v1beta.RegisterTraversalServiceHandler(ctx, gatewayMux, gatewayConn)

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

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
