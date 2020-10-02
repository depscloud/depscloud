package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/depscloud/api/swagger"
	"github.com/depscloud/api/v1alpha/extractor"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/gateway/internal/checks"
	"github.com/depscloud/depscloud/gateway/internal/proxies"
	"github.com/depscloud/depscloud/internal/mux"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/health"
)

// variables set during build using -X ldflag
var version string
var commit string
var date string

// https://github.com/grpc/grpc/blob/master/doc/service_config.md
const serviceConfigTemplate = `{
	"loadBalancingPolicy": "%s",
	"healthCheckConfig": {
		"serviceName": ""
	}
}`

func dial(target, certFile, keyFile, caFile, lbPolicy string) (*grpc.ClientConn, error) {
	serviceConfig := fmt.Sprintf(serviceConfigTemplate, lbPolicy)

	dialOptions := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(serviceConfig),
	}

	if len(certFile) > 0 {
		certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		certPool := x509.NewCertPool()
		bs, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		ok := certPool.AppendCertsFromPEM(bs)
		if !ok {
			return nil, fmt.Errorf("failed to append certs")
		}

		transportCreds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})

		dialOptions = append(dialOptions, grpc.WithTransportCredentials(transportCreds))
	} else {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	return grpc.Dial(target, dialOptions...)
}

func dialExtractor(cfg *gatewayConfig) (*grpc.ClientConn, error) {
	return dial(cfg.extractorAddress,
		cfg.extractorCertPath, cfg.extractorKeyPath, cfg.extractorCAPath,
		cfg.extractorLBPolicy)
}

func dialTracker(cfg *gatewayConfig) (*grpc.ClientConn, error) {
	return dial(cfg.trackerAddress,
		cfg.trackerCertPath, cfg.trackerKeyPath, cfg.trackerCAPath,
		cfg.trackerLBPolicy)
}

type gatewayConfig struct {
	httpPort int
	grpcPort int

	extractorAddress  string
	extractorCertPath string
	extractorKeyPath  string
	extractorCAPath   string
	extractorLBPolicy string

	trackerAddress  string
	trackerCertPath string
	trackerKeyPath  string
	trackerCAPath   string
	trackerLBPolicy string
}

func main() {
	cfg := &gatewayConfig{
		httpPort: 8080,
		grpcPort: 8090,

		extractorAddress:  "extractor:8090",
		extractorCertPath: "",
		extractorKeyPath:  "",
		extractorCAPath:   "",
		extractorLBPolicy: "round_robin",

		trackerAddress:  "tracker:8090",
		trackerCertPath: "",
		trackerKeyPath:  "",
		trackerCAPath:   "",
		trackerLBPolicy: "round_robin",
	}

	tlsConfig := &mux.TLSConfig{}

	app := &cli.App{
		Name:  "gateway",
		Usage: "an HTTP/gRPC proxy to backend services",
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Output version information",
				Action: func(c *cli.Context) error {
					versionString := fmt.Sprintf("%s %s", c.Command.Name, mux.Version{Version: version, Commit: commit, Date: date})

					fmt.Println(versionString)
					return nil

				},
			},
		},
		Flags: []cli.Flag{
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
				Name:        "extractor-address",
				Usage:       "address to the extractor service",
				Value:       cfg.extractorAddress,
				Destination: &cfg.extractorAddress,
				EnvVars:     []string{"EXTRACTOR_ADDRESS"},
			},
			&cli.StringFlag{
				Name:        "extractor-cert",
				Usage:       "certificate used to enable TLS for the extractor",
				Value:       cfg.extractorCertPath,
				Destination: &cfg.extractorCertPath,
				EnvVars:     []string{"EXTRACTOR_CERT_PATH"},
			},
			&cli.StringFlag{
				Name:        "extractor-key",
				Usage:       "key used to enable TLS for the extractor",
				Value:       cfg.extractorKeyPath,
				Destination: &cfg.extractorKeyPath,
				EnvVars:     []string{"EXTRACTOR_KEY_PATH"},
			},
			&cli.StringFlag{
				Name:        "extractor-ca",
				Usage:       "ca used to enable TLS for the extractor",
				Value:       cfg.extractorCAPath,
				Destination: &cfg.extractorCAPath,
				EnvVars:     []string{"EXTRACTOR_CA_PATH"},
			},
			&cli.StringFlag{
				Name:        "extractor-lb",
				Usage:       "the load balancer policy to use for the extractor",
				Value:       cfg.extractorLBPolicy,
				Destination: &cfg.extractorLBPolicy,
				EnvVars:     []string{"EXTRACTOR_LBPOLICY"},
			},
			&cli.StringFlag{
				Name:        "tracker-address",
				Usage:       "address to the tracker service",
				Value:       cfg.trackerAddress,
				Destination: &cfg.trackerAddress,
				EnvVars:     []string{"TRACKER_ADDRESS"},
			},
			&cli.StringFlag{
				Name:        "tracker-cert",
				Usage:       "certificate used to enable TLS for the tracker",
				Value:       cfg.trackerCertPath,
				Destination: &cfg.trackerCertPath,
				EnvVars:     []string{"TRACKER_CERT_PATH"},
			},
			&cli.StringFlag{
				Name:        "tracker-key",
				Usage:       "key used to enable TLS for the tracker",
				Value:       cfg.trackerKeyPath,
				Destination: &cfg.trackerKeyPath,
				EnvVars:     []string{"TRACKER_KEY_PATH"},
			},
			&cli.StringFlag{
				Name:        "tracker-ca",
				Usage:       "ca used to enable TLS for the tracker",
				Value:       cfg.trackerCAPath,
				Destination: &cfg.trackerCAPath,
				EnvVars:     []string{"TRACKER_CA_PATH"},
			},
			&cli.StringFlag{
				Name:        "tracker-lb",
				Usage:       "the load balancer policy to use for the tracker",
				Value:       cfg.trackerLBPolicy,
				Destination: &cfg.trackerLBPolicy,
				EnvVars:     []string{"TRACKER_LBPOLICY"},
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
		},
		Action: func(c *cli.Context) error {
			grpcServer, httpServer := mux.DefaultServers()
			gatewayMux := runtime.NewServeMux()

			ctx := context.Background()

			extractorConn, err := dialExtractor(cfg)
			if err != nil {
				return err
			}
			defer extractorConn.Close()

			trackerConn, err := dialTracker(cfg)
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
				Version:         &mux.Version{Version: version, Commit: commit, Date: date},
				TLSConfig:       tlsConfig,
			})
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
