package main

import (
	"fmt"
	"github.com/deps-cloud/api/swagger"
	"net/http"
	"os"
	"strings"

	"github.com/deps-cloud/api/v1alpha/extractor"
	"github.com/deps-cloud/api/v1alpha/tracker"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/rs/cors"
)

func exitIff(err error) {
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
}

func dialOptions(cert string) []grpc.DialOption {
	opts := make([]grpc.DialOption, 0)
	if len(cert) > 0 {
		transportCreds, err := credentials.NewClientTLSFromFile(cert, "")
		exitIff(err)

		opts = append(opts, grpc.WithTransportCredentials(transportCreds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return opts
}

func prefixedHandle(prefix string, mux http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.URL.Path = strings.TrimPrefix(request.URL.Path, prefix)
		mux.ServeHTTP(writer, request)
	}
}

func

func main() {
	port := 8080
	extractorAddress := "extractor:8090"
	extractorCert := ""
	trackerAddress := "tracker:8090"
	trackerCert := ""
	tlsCert := ""
	tlsKey := ""

	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "Start up an HTTP proxy for the gRPC services",
		Run: func(cmd *cobra.Command, args []string) {
			address := fmt.Sprintf(":%d", port)

			swaggerMux := http.FileServer(swagger.AssetFile())
			gatewayMux := runtime.NewServeMux()

			ctx := context.Background()

			trackerOpts := dialOptions(trackerCert)
			extractorOpts := dialOptions(extractorCert)

			err := tracker.RegisterSourceServiceHandlerFromEndpoint(ctx, gatewayMux, trackerAddress, trackerOpts)
			exitIff(err)

			err = tracker.RegisterModuleServiceHandlerFromEndpoint(ctx, gatewayMux, trackerAddress, trackerOpts)
			exitIff(err)

			err = tracker.RegisterDependencyServiceHandlerFromEndpoint(ctx, gatewayMux, trackerAddress, trackerOpts)
			exitIff(err)

			err = tracker.RegisterTopologyServiceHandlerFromEndpoint(ctx, gatewayMux, trackerAddress, trackerOpts)
			exitIff(err)

			err = extractor.RegisterDependencyExtractorHandlerFromEndpoint(ctx, gatewayMux, extractorAddress, extractorOpts)
			exitIff(err)

			httpMux := http.NewServeMux()

			httpMux.HandleFunc("/swagger/", prefixedHandle("/swagger", swaggerMux))
			httpMux.Handle("/", gatewayMux)

			apiMux := cors.Default().Handler(httpMux)

			if len(tlsCert) > 0 && len(tlsKey) > 0 {
				logrus.Infof("[main] starting TLS server on %s", address)
				err = http.ListenAndServeTLS(address, tlsCert, tlsKey, apiMux)
			} else {
				logrus.Infof("[main] starting plaintext server on %s", address)
				err = http.ListenAndServe(address, apiMux)
			}
			exitIff(err)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&port, "port", port, "(optional) the port to run on")
	flags.StringVar(&extractorAddress, "extractor-address", extractorAddress, "(optional) address to des")
	flags.StringVar(&extractorCert, "extractor-cert", extractorCert, "(optional) certificate used to enable TLS for the extractor")
	flags.StringVar(&trackerAddress, "tracker-address", trackerAddress, "(optional) address to dts")
	flags.StringVar(&trackerCert, "tracker-cert", trackerCert, "(optional) certificate used to enable TLS for the tracker")
	flags.StringVar(&tlsKey, "tls-key", tlsKey, "(optional) path to the file containing the TLS private key")
	flags.StringVar(&tlsCert, "tls-cert", tlsCert, "(optional) path to the file containing the TLS certificate")

	err := cmd.Execute()
	exitIff(err)
}
