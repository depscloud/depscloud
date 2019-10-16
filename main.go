package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/deps-cloud/api/v1alpha/extractor"
	"github.com/deps-cloud/api/v1alpha/tracker"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

			mux := runtime.NewServeMux()

			ctx := context.Background()

			trackerOpts := dialOptions(trackerCert)
			extractorOpts := dialOptions(extractorCert)

			err := tracker.RegisterSourceServiceHandlerFromEndpoint(ctx, mux, trackerAddress, trackerOpts)
			exitIff(err)

			err = tracker.RegisterModuleServiceHandlerFromEndpoint(ctx, mux, trackerAddress, trackerOpts)
			exitIff(err)

			err = tracker.RegisterDependencyServiceHandlerFromEndpoint(ctx, mux, trackerAddress, trackerOpts)
			exitIff(err)

			err = tracker.RegisterTopologyServiceHandlerFromEndpoint(ctx, mux, trackerAddress, trackerOpts)
			exitIff(err)

			err = extractor.RegisterDependencyExtractorHandlerFromEndpoint(ctx, mux, extractorAddress, extractorOpts)
			exitIff(err)

			if len(tlsCert) > 0 && len(tlsKey) > 0 {
				logrus.Infof("[main] starting TLS server on %s", address)
				err = http.ListenAndServeTLS(address, tlsCert, tlsKey, mux)
			} else {
				logrus.Infof("[main] starting plaintext server on %s", address)
				err = http.ListenAndServe(address, mux)
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
