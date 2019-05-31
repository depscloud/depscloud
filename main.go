package main

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/pkg/service"
	"github.com/deps-cloud/dts/pkg/store"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func panicIff(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	configPath := "${HOME}/.dts/config.yaml"
	port := 8090

	cmd := &cobra.Command{
		Use: "",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			address := fmt.Sprintf(":%d", port)

			listener, err := net.Listen("tcp", address)
			panicIff(err)

			// todo: make this configurable
			db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
			panicIff(err)

			graphStore, err := store.NewSQLGraphStore(db)
			panicIff(err)

			dts, err := service.NewDependencyTrackingService(graphStore)
			panicIff(err)

			server := grpc.NewServer()
			api.RegisterDependencyTrackingServiceServer(server, dts)

			logrus.Infof("[main] starting gRPC on %s", address)
			err = server.Serve(listener)
			panicIff(err)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&configPath, "config", configPath, "The path to the config file")
	flags.IntVar(&port, "port", port, "The port to run on")

	err := cmd.Execute()
	panicIff(err)
}
