package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/depscloud/api/v1alpha/store"
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
	"google.golang.org/grpc/credentials"
)

func registerV1Alpha(graphStoreClient store.GraphStoreClient, server *grpc.Server) {
	// v1alpha
	services.RegisterDependencyService(server, graphStoreClient)
	services.RegisterModuleService(server, graphStoreClient)
	services.RegisterSourceService(server, graphStoreClient)
	services.RegisterSearchService(server, graphStoreClient)
}

type trackerConfig struct {
	port                   int
	storageDriver          string
	storageAddress         string
	storageReadOnlyAddress string

	tlsKeyPath  string
	tlsCertPath string
	tlsCAPath   string
}

var description = strings.TrimSpace(`
   To learn more about how to configure the storage layer, see our documentation.
   https://deps.cloud/docs/deploy/config/storage/
`)

func main() {
	cfg := &trackerConfig{
		port:                   8090,
		storageDriver:          "sqlite",
		storageAddress:         "file::memory:?cache=shared",
		storageReadOnlyAddress: "",
		tlsKeyPath:             "",
		tlsCertPath:            "",
		tlsCAPath:              "",
	}

	app := &cli.App{
		Name:        "tracker",
		Usage:       "tracks dependencies between systems",
		Description: description,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Usage:       "the port to run on",
				Value:       cfg.port,
				Destination: &cfg.port,
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
				Value:       cfg.tlsKeyPath,
				Destination: &cfg.tlsKeyPath,
				EnvVars:     []string{"TLS_KEY_PATH"},
			},
			&cli.StringFlag{
				Name:        "tls-cert",
				Usage:       "path to the file containing the TLS certificate",
				Value:       cfg.tlsCertPath,
				Destination: &cfg.tlsCertPath,
				EnvVars:     []string{"TLS_CERT_PATH"},
			},
			&cli.StringFlag{
				Name:        "tls-ca",
				Usage:       "path to the file containing the TLS certificate authority",
				Value:       cfg.tlsCAPath,
				Destination: &cfg.tlsCAPath,
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

			options := make([]grpc.ServerOption, 0)
			if len(cfg.tlsCertPath) > 0 && len(cfg.tlsKeyPath) > 0 && len(cfg.tlsCAPath) > 0 {
				logrus.Info("[main] configuring tls")

				certificate, err := tls.LoadX509KeyPair(cfg.tlsCertPath, cfg.tlsKeyPath)
				if err != nil {
					return err
				}

				certPool := x509.NewCertPool()
				bs, err := ioutil.ReadFile(cfg.tlsCAPath)
				if err != nil {
					return err
				}

				ok := certPool.AppendCertsFromPEM(bs)
				if !ok {
					return fmt.Errorf("failed to append certs")
				}

				transportCreds := credentials.NewTLS(&tls.Config{
					ClientAuth:   tls.RequireAndVerifyClientCert,
					Certificates: []tls.Certificate{certificate},
					ClientCAs:    certPool,
				})

				options = append(options, grpc.Creds(transportCreds))
			}

			graphStore, err := graphstore.NewSQLGraphStore(rwdb, rodb, statements)
			if err != nil {
				return err
			}

			graphStoreClient := store.NewInProcessGraphStoreClient(graphStore)
			graphStoreClient = graphstore.Retryable(graphStoreClient, 5)

			grpcServer := grpc.NewServer(options...)
			registerV1Alpha(graphStoreClient, grpcServer)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			allChecks := checks.Checks(graphStoreClient)
			checks.RegisterHealthCheck(ctx, grpcServer, allChecks)

			// setup server
			address := fmt.Sprintf(":%d", cfg.port)

			listener, err := net.Listen("tcp", address)
			if err != nil {
				return err
			}

			logrus.Infof("[main] starting gRPC on %s", address)
			return grpcServer.Serve(listener)
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
