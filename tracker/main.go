package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	apiv1alpha "github.com/depscloud/api/v1alpha/store"
	apiv1beta "github.com/depscloud/api/v1beta/graphstore"
	"github.com/depscloud/depscloud/internal/logger"
	"github.com/depscloud/depscloud/internal/mux"
	"github.com/depscloud/depscloud/internal/v"
	"github.com/depscloud/depscloud/tracker/internal/checks"
	"github.com/depscloud/depscloud/tracker/internal/cleanup"
	"github.com/depscloud/depscloud/tracker/internal/graphstore/v1alpha"
	"github.com/depscloud/depscloud/tracker/internal/graphstore/v1beta"
	svcsv1alpha "github.com/depscloud/depscloud/tracker/internal/services/v1alpha"
	svcsv1beta "github.com/depscloud/depscloud/tracker/internal/services/v1beta"

	_ "github.com/go-sql-driver/mysql"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	_ "github.com/mattn/go-sqlite3"

	"github.com/urfave/cli/v2"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

// variables set during build using -X ldflag
var version string
var commit string
var date string

const sockAddr = "localhost:47274"

func graphStoreServers(driver, address, readOnlyAddress string) (apiv1alpha.GraphStoreServer, apiv1beta.GraphStoreServer, error) {
	// v1beta
	v1betaDriver, err := v1beta.Resolve(driver, address, readOnlyAddress)
	if err != nil {
		return nil, nil, err
	}

	// v1alpha
	v1alphaGraphStore, err := v1alpha.NewGraphStoreFor(driver, address, readOnlyAddress)
	if err != nil {
		return nil, nil, err
	}

	return v1alphaGraphStore, &v1beta.GraphStoreServer{Driver: v1betaDriver}, nil
}

func startGraphStore(log *zap.Logger, driver, address, readOnlyAddress string) error {
	grpcServer := grpc.NewServer()

	v1alphaGraphStore, v1betaGraphStore, err := graphStoreServers(driver, address, readOnlyAddress)
	if err != nil {
		return err
	}

	apiv1beta.RegisterGraphStoreServer(grpcServer, v1betaGraphStore)
	apiv1alpha.RegisterGraphStoreServer(grpcServer, v1alphaGraphStore)

	// listen and serve
	log.Info("starting graphstore grpc", zap.String("bind", sockAddr))
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
	svcsv1beta.RegisterManifestStorageServiceServer(server, v1betaClient)
	// more eventually
}

type trackerConfig struct {
	storageDriver          string
	storageAddress         string
	storageReadOnlyAddress string
}

var description = strings.TrimSpace(`
   To learn more about how to configure the storage layer, see our documentation.
   https://deps.cloud/docs/deploy/config/storage/
`)

func main() {
	version := v.Info{Version: version, Commit: commit, Date: date}

	loggerConfig, loggerFlags := logger.WithFlags(zap.NewProductionConfig())
	serverConfig, serverFlags := mux.WithFlags(mux.DefaultConfig(version))

	cfg := &trackerConfig{
		storageDriver:          "sqlite",
		storageAddress:         "file::memory:?cache=shared",
		storageReadOnlyAddress: "",
	}

	flags := []cli.Flag{
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
	}

	flags = append(flags, loggerFlags...)
	flags = append(flags, serverFlags...)

	app := &cli.App{
		Name:        "tracker",
		Usage:       "tracks dependencies between systems",
		Description: description,
		Commands: []*cli.Command{
			{
				Name:  "cleanup",
				Usage: "Cleanup data in the database",
				Flags: flags,
				Action: func(context *cli.Context) error {
					v1alphaGraphStore, v1betaGraphStore, err := graphStoreServers(
						cfg.storageDriver, cfg.storageAddress, cfg.storageReadOnlyAddress)

					if err != nil {
						return err
					}

					servers := cleanup.NewServers(v1alphaGraphStore, v1betaGraphStore)
					return cleanup.Run(servers)
				},
			},
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

			err := startGraphStore(log, cfg.storageDriver, cfg.storageAddress, cfg.storageReadOnlyAddress)
			if err != nil {
				return err
			}

			cc, err := grpc.Dial(sockAddr,
				grpc.WithInsecure(),
				grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
					grpc_prometheus.StreamClientInterceptor,
				)),
				grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
					grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(5)),
					grpc_prometheus.UnaryClientInterceptor,
				)),
			)
			if err != nil {
				return err
			}

			v1alphaClient := apiv1alpha.NewGraphStoreClient(cc)
			v1betaClient := apiv1beta.NewGraphStoreClient(cc)

			// setup checks and any extra endpoints
			serverConfig.Checks = checks.Checks(v1betaClient, v1alphaClient)
			serverConfig.Endpoints = []mux.ServerEndpoint{
				func(ctx context.Context, grpcServer *grpc.Server, httpServer *http.ServeMux) {
					registerV1Alpha(v1alphaClient, grpcServer)
					registerV1Beta(v1betaClient, grpcServer)
				},
			}

			server := mux.NewServer(serverConfig)
			return server.Serve(ctx)
		},
	}

	_ = app.Run(os.Args)
}
