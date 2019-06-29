package main

import (
	"fmt"
	"net/http"

	desapi "github.com/deps-cloud/des/api"
	dtsapi "github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/api/v1alpha"
	rdsapi "github.com/deps-cloud/rds/api"

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
	rdsAddress := "rds:8090"
	desAddress := "des:8090"
	dtsAddress := "dts:8090"

	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			address := fmt.Sprintf(":%d", port)

			opts := []grpc.DialOption{
				grpc.WithInsecure(),
			}

			mux := runtime.NewServeMux()

			ctx := context.Background()

			err := v1alpha.RegisterSourceServiceHandlerFromEndpoint(ctx, mux, dtsAddress, opts)
			panicIff(err)

			err = v1alpha.RegisterModuleServiceHandlerFromEndpoint(ctx, mux, dtsAddress, opts)
			panicIff(err)

			err = v1alpha.RegisterDependencyServiceHandlerFromEndpoint(ctx, mux, dtsAddress, opts)
			panicIff(err)

			err = v1alpha.RegisterTopologyServiceHandlerFromEndpoint(ctx, mux, dtsAddress, opts)
			panicIff(err)

			err = dtsapi.RegisterDependencyTrackerHandlerFromEndpoint(ctx, mux, dtsAddress, opts)
			panicIff(err)

			err = desapi.RegisterDependencyExtractorHandlerFromEndpoint(ctx, mux, desAddress, opts)
			panicIff(err)

			err = rdsapi.RegisterRepositoryDiscoveryHandlerFromEndpoint(ctx, mux, rdsAddress, opts)
			panicIff(err)

			err = http.ListenAndServe(address, mux)
			panicIff(err)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&port, "port", port, "(optional) the port to run on")
	flags.StringVar(&rdsAddress, "rds-address", rdsAddress, "(optional) address to rds")
	flags.StringVar(&desAddress, "des-address", desAddress, "(optional) address to des")
	flags.StringVar(&dtsAddress, "dts-address", dtsAddress, "(optional) address to dts")

	err := cmd.Execute()
	panicIff(err)
}
