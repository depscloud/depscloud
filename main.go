package main

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/api/v1alpha/store"
	"github.com/deps-cloud/dts/pkg/service"
	"github.com/deps-cloud/dts/pkg/services"
	"github.com/deps-cloud/dts/pkg/services/graphstore"

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

func registerV1Alpha(rwdb, rodb *sql.DB, server *grpc.Server) {
	graphStore, err := graphstore.NewSQLGraphStore(rwdb, rodb)
	panicIff(err)

	graphStoreClient := store.NewInProcessGraphStoreClient(graphStore)

	// poc
	dts, _ := service.NewDependencyTrackingService(graphStoreClient)
	api.RegisterDependencyTrackerServer(server, dts)

	// v1alpha
	services.RegisterDependencyService(server, graphStoreClient)
	services.RegisterModuleService(server, graphStoreClient)
	services.RegisterSourceService(server, graphStoreClient)
	services.RegisterTopologyService(server, graphStoreClient)
}

func main() {
	configPath := "${HOME}/.dts/config.yaml"
	port := 8090
	storageDriver := "sqlite3"
	storageAddress := "file::memory:?cache=shared"
	storageReadOnlyAddress := ""

	cmd := &cobra.Command{
		Use:   "dts",
		Short: "dts runs the dependency tracking service.",
		Run: func(cmd *cobra.Command, args []string) {
			rwdb, err := sql.Open(storageDriver, storageAddress)
			panicIff(err)

			rodb := rwdb
			if len(storageReadOnlyAddress) > 0 {
				rodb, err = sql.Open(storageDriver, storageReadOnlyAddress)
				panicIff(err)
			}

			server := grpc.NewServer()
			healthpb.RegisterHealthServer(server, health.NewServer())
			registerV1Alpha(rwdb, rodb, server)

			// setup server
			address := fmt.Sprintf(":%d", port)

			listener, err := net.Listen("tcp", address)
			panicIff(err)

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
	flags.StringVar(&storageReadOnlyAddress, "storage-readonly-address", storageReadOnlyAddress, "(optional) the readonly address of the storage tier")

	err := cmd.Execute()
	panicIff(err)
}
