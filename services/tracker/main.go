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
	"github.com/depscloud/depscloud/internal/appconf"
	"github.com/depscloud/depscloud/internal/logger"
	"github.com/depscloud/depscloud/internal/mux"
	"github.com/depscloud/depscloud/services/tracker/internal/checks"
	"github.com/depscloud/depscloud/services/tracker/internal/cleanup"
	"github.com/depscloud/depscloud/services/tracker/internal/db"
	"github.com/depscloud/depscloud/services/tracker/internal/graphstore/v1alpha"
	"github.com/depscloud/depscloud/services/tracker/internal/graphstore/v1beta"
	svcsv1alpha "github.com/depscloud/depscloud/services/tracker/internal/services/v1alpha"
	svcsv1beta "github.com/depscloud/depscloud/services/tracker/internal/services/v1beta"

	_ "github.com/go-sql-driver/mysql"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	_ "github.com/mattn/go-sqlite3"

	"github.com/urfave/cli/v2"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"

	"gorm.io/gorm"
)

const sockAddr = "0.0.0.0:47274"

func graphStoreServers(name string, rw, ro *gorm.DB) (apiv1alpha.GraphStoreServer, apiv1beta.GraphStoreServer, error) {
	v1alphaStatements := db.StatementsFor(name, "v1alpha")
	v1betaStatements := db.StatementsFor(name, "v1beta")

	if rw != nil {
		err := rw.AutoMigrate(
			&v1beta.GraphData{},
		)

		if err != nil {
			return nil, nil, err
		}
	}

	sqlxRW, err := db.ToSQLX(name, rw)
	if err != nil {
		return nil, nil, err
	}

	sqlxRO, err := db.ToSQLX(name, ro)
	if err != nil {
		return nil, nil, err
	}

	v1alphaGraphStore, err := v1alpha.NewSQLGraphStore(sqlxRW, sqlxRO, v1alphaStatements)
	if err != nil {
		return nil, nil, err
	}

	v1betaGraphStore := &v1beta.GraphStoreServer{
		Driver: v1beta.NewSQLDriver(sqlxRW, sqlxRO, v1betaStatements),
	}

	return v1alphaGraphStore, v1betaGraphStore, nil
}

func startGraphStore(log *zap.Logger, driver string, rw, ro *gorm.DB) error {
	grpcServer := grpc.NewServer()

	v1alphaGraphStore, v1betaGraphStore, err := graphStoreServers(driver, rw, ro)
	if err != nil {
		return err
	}

	apiv1beta.RegisterGraphStoreServer(grpcServer, v1betaGraphStore)
	apiv1alpha.RegisterGraphStoreServer(grpcServer, v1alphaGraphStore)

	log.Info("starting graphstore",
		zap.String("bind", sockAddr),
		zap.String("protocol", "grpc"))

	listener, err := net.Listen("tcp", sockAddr)
	if err != nil {
		return err
	}

	go grpcServer.Serve(listener)
	return nil
}

func registerV1Alpha(server *grpc.Server, v1alphaClient apiv1alpha.GraphStoreClient) {
	svcsv1alpha.RegisterDependencyService(server, v1alphaClient)
	svcsv1alpha.RegisterModuleService(server, v1alphaClient)
	svcsv1alpha.RegisterSourceService(server, v1alphaClient)
	svcsv1alpha.RegisterSearchService(server, v1alphaClient)
}

func registerV1Beta(server *grpc.Server, v1betaClient apiv1beta.GraphStoreClient) {
	svcsv1beta.RegisterManifestStorageServiceServer(server, v1betaClient)
	svcsv1beta.RegisterModuleServiceServer(server, v1betaClient)
	svcsv1beta.RegisterSourceServiceServer(server, v1betaClient)
	svcsv1beta.RegisterTraversalServiceServer(server, v1betaClient)
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
	version := appconf.Current()

	loggerConfig, loggerFlags := logger.WithFlags(logger.DefaultConfig())
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
					name, rw, ro, err := db.Resolve(cfg.storageDriver, cfg.storageAddress, cfg.storageReadOnlyAddress)
					if err != nil {
						return err
					}

					v1alphaGraphStore, v1betaGraphStore, err := graphStoreServers(name, rw, ro)
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

			name, rw, ro, err := db.Resolve(cfg.storageDriver, cfg.storageAddress, cfg.storageReadOnlyAddress)
			if err != nil {
				return err
			}

			err = startGraphStore(log, name, rw, ro)
			if err != nil {
				return err
			}

			cc, err := grpc.Dial(sockAddr,
				grpc.WithInsecure(),
				grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
					grpc_prometheus.StreamClientInterceptor,
					grpc_zap.StreamClientInterceptor(log),
				)),
				grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
					grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(5)),
					grpc_prometheus.UnaryClientInterceptor,
					grpc_zap.UnaryClientInterceptor(log),
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
					registerV1Alpha(grpcServer, v1alphaClient)
					registerV1Beta(grpcServer, v1betaClient)
				},
			}

			server := mux.NewServer(serverConfig)
			return server.Serve(ctx)
		},
	}

	_ = app.Run(os.Args)
}
