package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	apiv1alpha "github.com/depscloud/api/v1alpha/store"
	apiv1beta "github.com/depscloud/api/v1beta/graphstore"
	"github.com/depscloud/depscloud/internal/mux"
	"github.com/depscloud/depscloud/tracker/internal/checks"
	"github.com/depscloud/depscloud/tracker/internal/graphstore/v1alpha"
	"github.com/depscloud/depscloud/tracker/internal/graphstore/v1beta"
	svcsv1alpha "github.com/depscloud/depscloud/tracker/internal/services/v1alpha"

	_ "github.com/go-sql-driver/mysql"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"

	"google.golang.org/grpc"
)

// variables set during build using -X ldflag
var version string
var commit string
var date string

const sockAddr = "localhost:47274"

func startGraphStore(driver, address, readOnlyAddress string) error {
	grpcServer := grpc.NewServer()

	// v1beta
	v1betaDriver, err := v1beta.Resolve(driver, address, readOnlyAddress)
	if err != nil {
		return err
	}
	apiv1beta.RegisterGraphStoreServer(grpcServer, &v1beta.GraphStoreServer{Driver: v1betaDriver})

	// v1alpha
	v1alphaGraphStore, err := v1alpha.NewGraphStoreFor(driver, address, readOnlyAddress)
	if err != nil {
		return err
	}
	apiv1alpha.RegisterGraphStoreServer(grpcServer, v1alphaGraphStore)

	// listen and serve
	logrus.Infof("[graphstore] starting grpc on %s", sockAddr)
	listener, err := net.Listen("tcp", sockAddr)
	if err != nil {
		return err
	}

	go grpcServer.Serve(listener)
	return nil
}

func registerV1Alpha(v1alphaClient apiv1alpha.GraphStoreClient, server *grpc.Server) {
	svcsv1alpha.RegisterDependencyService(server, v1alphaClient)
	svcsv1alpha.RegisterModuleService(server, v1alphaClient)
	svcsv1alpha.RegisterSourceService(server, v1alphaClient)
	svcsv1alpha.RegisterSearchService(server, v1alphaClient)
}

func registerV1Beta(v1betaClient apiv1beta.GraphStoreClient, server *grpc.Server) {
	// TODO: fill in with gh-50 to gh-57
}

type trackerConfig struct {
	httpPort               int
	grpcPort               int
	storageDriver          string
	storageAddress         string
	storageReadOnlyAddress string
}

var description = strings.TrimSpace(`
   To learn more about how to configure the storage layer, see our documentation.
   https://deps.cloud/docs/deploy/config/storage/
`)

func main() {
	version := mux.Version{Version: version, Commit: commit, Date: date}
	cfg := &trackerConfig{
		httpPort:               8080,
		grpcPort:               8090,
		storageDriver:          "sqlite",
		storageAddress:         "file::memory:?cache=shared",
		storageReadOnlyAddress: "",
	}

	tlsConfig := &mux.TLSConfig{}

	app := &cli.App{
		Name:        "tracker",
		Usage:       "tracks dependencies between systems",
		Description: description,
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Output version information",
				Action: func(c *cli.Context) error {
					versionString := fmt.Sprintf("%s %s", c.Command.Name, version)
					fmt.Println(versionString)
					return nil

				},
			},
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "http-port",
				Usage:       "the port to run http on",
				Value:       cfg.httpPort,
				Destination: &cfg.httpPort,
				EnvVars:     []string{"HTTP_PORT"},
			},
			&cli.IntFlag{
				Name:        "grpc-port",
				Aliases:     []string{"port"},
				Usage:       "the port to run grpc on",
				Value:       cfg.grpcPort,
				Destination: &cfg.grpcPort,
				EnvVars:     []string{"GRPC_PORT"},
			},
			&cli.StringFlag{
				Name:        "storage-driver",
				Usage:       "the driver used to configure the storage tier",
				Value:       cfg.storageDriver,
				Destination: &cfg.storageDriver,
				EnvVars:     []string{"STORAGE_DRIVER"},
			},
			&cli.StringFlag{
				Name:        "storage-address",
				Usage:       "the address of the storage tier",
				Value:       cfg.storageAddress,
				Destination: &cfg.storageAddress,
				EnvVars:     []string{"STORAGE_ADDRESS"},
			},
			&cli.StringFlag{
				Name:        "storage-readonly-address",
				Usage:       "the readonly address of the storage tier",
				Value:       cfg.storageReadOnlyAddress,
				Destination: &cfg.storageReadOnlyAddress,
				EnvVars:     []string{"STORAGE_READ_ONLY_ADDRESS"},
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
			err := startGraphStore(cfg.storageDriver, cfg.storageAddress, cfg.storageReadOnlyAddress)
			if err != nil {
				return err
			}

			cc, err := grpc.Dial(sockAddr,
				grpc.WithInsecure(),
				grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
					grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(5)),
					grpc_prometheus.UnaryClientInterceptor,
				)),
				grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
					grpc_prometheus.StreamClientInterceptor,
				)),
			)
			if err != nil {
				return err
			}

			grpcServer, httpServer := mux.DefaultServers()

			v1betaClient := apiv1beta.NewGraphStoreClient(cc)
			registerV1Beta(v1betaClient, grpcServer)

			v1alphaClient := apiv1alpha.NewGraphStoreClient(cc)
			registerV1Alpha(v1alphaClient, grpcServer)

			return mux.Serve(grpcServer, httpServer, &mux.Config{
				Context:         c.Context,
				BindAddressHTTP: fmt.Sprintf("0.0.0.0:%d", cfg.httpPort),
				BindAddressGRPC: fmt.Sprintf("0.0.0.0:%d", cfg.grpcPort),
				Checks:          checks.Checks(v1betaClient, v1alphaClient),
				Version:         &version,
				TLSConfig:       tlsConfig,
			})
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
