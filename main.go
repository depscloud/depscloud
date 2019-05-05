package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
)

func main() {
	configPath := "${HOME}/.dts/config.yaml"
	port := 8090

	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			address := fmt.Sprintf(":%d", port)

			listener, err := net.Listen("tcp", address)
			if err != nil {
				panic(err)
			}

			server := grpc.NewServer()
			//api.RegisterDependencyTrackingServiceServer(server, impl)

			logrus.Infof("[main] starting gRPC on %s", address)
			if err := server.Serve(listener); err != nil {
				panic(err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&configPath, "config", configPath, "The path to the config file")
	flags.IntVar(&port, "port", port, "The port to run on")

	if err := cmd.Execute(); err != nil {
		logrus.Errorf("")
		panic(err)
	}
}
