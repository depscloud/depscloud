package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/depscloud/internal/mux"
	"github.com/depscloud/depscloud/tracker/internal/checks"
	"github.com/depscloud/depscloud/tracker/internal/graphstore"
	"github.com/depscloud/depscloud/tracker/internal/services"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"

	"google.golang.org/grpc"
)

// variables set during build using -X ldflag
var version string
var commit string
var date string

func registerV1Alpha(graphStoreClient store.GraphStoreClient, server *grpc.Server) {
	// v1alpha
	services.RegisterDependencyService(server, graphStoreClient)
	services.RegisterModuleService(server, graphStoreClient)
	services.RegisterSourceService(server, graphStoreClient)
	services.RegisterSearchService(server, graphStoreClient)
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
					versionString := fmt.Sprintf("{version: %s, commit: %s, date: %s}", version, commit, date)
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
			var rwdb *sqlx.DB
			var err error

			cfg.storageDriver, err = graphstore.ResolveDriverName(cfg.storageDriver)
			if err != nil {
				return err
			}

			if len(cfg.storageAddress) > 0 {
				rwdb, err = sqlx.Open(cfg.storageDriver, cfg.storageAddress)
				if err != nil {
					return err
				}
			}

			rodb := rwdb
			if len(cfg.storageReadOnlyAddress) > 0 {
				rodb, err = sqlx.Open(cfg.storageDriver, cfg.storageReadOnlyAddress)
				if err != nil {
					return err
				}
			}

			if rodb == nil && rwdb == nil {
				return fmt.Errorf("either --storage-address or --storage-readonly-address must be provided")
			}

			statements, err := graphstore.DefaultStatementsFor(cfg.storageDriver)
			if err != nil {
				return err
			}

			graphStore, err := graphstore.NewSQLGraphStore(rwdb, rodb, statements)
			if err != nil {
				return err
			}

			graphStoreClient := store.NewInProcessGraphStoreClient(graphStore)
			graphStoreClient = graphstore.Retryable(graphStoreClient, 5)

			grpcServer, httpServer := mux.DefaultServers()
			registerV1Alpha(graphStoreClient, grpcServer)

			return mux.Serve(grpcServer, httpServer, &mux.Config{
				Context:         c.Context,
				BindAddressHTTP: fmt.Sprintf("0.0.0.0:%d", cfg.httpPort),
				BindAddressGRPC: fmt.Sprintf("0.0.0.0:%d", cfg.grpcPort),
				Checks:          checks.Checks(graphStoreClient),
				Version:         &mux.Version{Version: version, Commit: commit, Date: date},
				TLSConfig:       tlsConfig,
			})
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
