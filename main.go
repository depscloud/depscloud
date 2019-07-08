package main

import (
	"fmt"
	"net/http"

	rdsapi "github.com/deps-cloud/discovery/api"
	desapi "github.com/deps-cloud/extractor/api"
	dtsapi "github.com/deps-cloud/tracker/api"
	"github.com/deps-cloud/tracker/api/v1alpha"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/spf13/cobra"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func panicIff(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	port := 8080
	discoveryAddress := "discovery:8090"
	extractorAddress := "extractor:8090"
	trackerAddress := "tracker:8090"

	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "Start up an HTTP proxy for the gRPC services",
		Run: func(cmd *cobra.Command, args []string) {
			address := fmt.Sprintf(":%d", port)

			opts := []grpc.DialOption{
				grpc.WithInsecure(),
			}

			mux := runtime.NewServeMux()

			ctx := context.Background()

			err := v1alpha.RegisterSourceServiceHandlerFromEndpoint(ctx, mux, trackerAddress, opts)
			panicIff(err)

			err = v1alpha.RegisterModuleServiceHandlerFromEndpoint(ctx, mux, trackerAddress, opts)
			panicIff(err)

			err = v1alpha.RegisterDependencyServiceHandlerFromEndpoint(ctx, mux, trackerAddress, opts)
			panicIff(err)

			err = v1alpha.RegisterTopologyServiceHandlerFromEndpoint(ctx, mux, trackerAddress, opts)
			panicIff(err)

			err = dtsapi.RegisterDependencyTrackerHandlerFromEndpoint(ctx, mux, trackerAddress, opts)
			panicIff(err)

			err = desapi.RegisterDependencyExtractorHandlerFromEndpoint(ctx, mux, extractorAddress, opts)
			panicIff(err)

			err = rdsapi.RegisterRepositoryDiscoveryHandlerFromEndpoint(ctx, mux, discoveryAddress, opts)
			panicIff(err)

			err = http.ListenAndServe(address, mux)
			panicIff(err)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&port, "port", port, "(optional) the port to run on")
	flags.StringVar(&discoveryAddress, "discovery-address", discoveryAddress, "(optional) address to rds")
	flags.StringVar(&extractorAddress, "extractor-address", extractorAddress, "(optional) address to des")
	flags.StringVar(&trackerAddress, "tracker-address", trackerAddress, "(optional) address to dts")

	err := cmd.Execute()
	panicIff(err)
}
