package main

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/pkg/service"
	"github.com/deps-cloud/dts/pkg/store"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func panicIff(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	configPath := "${HOME}/.dts/config.yaml"
	port := 8090
	storageDriver := "sqlite3"
	storageAddress := "file::memory:?cache=shared"

	cmd := &cobra.Command{
		Use: "dts",
		Short: "dts runs the dependency tracking service.",
		Run: func(cmd *cobra.Command, args []string) {
			address := fmt.Sprintf(":%d", port)

			listener, err := net.Listen("tcp", address)
			panicIff(err)

			db, err := sql.Open(storageDriver, storageAddress)
			panicIff(err)

			graphStore, err := store.NewSQLGraphStore(db)
			panicIff(err)

			dts, err := service.NewDependencyTrackingService(graphStore)
			panicIff(err)

			healthcheck := health.NewServer()
			// toggle the service health as such
			// healthcheck.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

			server := grpc.NewServer()
			api.RegisterDependencyTrackerServer(server, dts)
			healthpb.RegisterHealthServer(server, healthcheck)

			logrus.Infof("[main] starting gRPC on %s", address)
			err = server.Serve(listener)
			panicIff(err)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&configPath, "config", configPath, "(optional) the path to the config file")
	flags.IntVar(&port, "port", port, "(optional) the port to run on")
	flags.StringVar(&storageDriver, "storage-driver", storageDriver, "(optional) the driver used to configure the storage tier")
	flags.StringVar(&storageAddress, "storage-address", storageAddress, "(optional) the address of the storage tier")

	err := cmd.Execute()
	panicIff(err)
}
