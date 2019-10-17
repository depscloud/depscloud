package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/deps-cloud/api/swagger"
	"io/ioutil"
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

func dialOptions(certFile, keyFile, caFile string) []grpc.DialOption {
	opts := make([]grpc.DialOption, 0)
	if len(certFile) > 0 {
		certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
		exitIff(err)

		certPool := x509.NewCertPool()
		bs, err := ioutil.ReadFile(caFile)
		exitIff(err)

		ok := certPool.AppendCertsFromPEM(bs)
		if !ok {
			exitIff(fmt.Errorf("failed to append certs"))
		}

		transportCreds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{certificate},
			RootCAs: certPool,
		})

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

func main() {
	port := 8080

	extractorAddress := "extractor:8090"
	extractorCert := ""
	extractorKey := ""
	extractorCA := ""

	trackerAddress := "tracker:8090"
	trackerCert := ""
	trackerKey := ""
	trackerCA := ""

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

			trackerOpts := dialOptions(trackerCert, trackerKey, trackerCA)
			extractorOpts := dialOptions(extractorCert, extractorKey, extractorCA)

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

	flags.StringVar(&extractorAddress, "extractor-address", extractorAddress, "(optional) address to the extractor service")
	flags.StringVar(&extractorCert, "extractor-cert", extractorCert, "(optional) certificate used to enable TLS for the extractor")
	flags.StringVar(&extractorKey, "extractor-key", extractorKey, "(optional) key used to enable TLS for the extractor")
	flags.StringVar(&extractorCA, "extractor-ca", extractorCA, "(optional) ca used to enable TLS for the extractor")

	flags.StringVar(&trackerAddress, "tracker-address", trackerAddress, "(optional) address to the tracker service")
	flags.StringVar(&trackerCert, "tracker-cert", trackerCert, "(optional) certificate used to enable TLS for the tracker")
	flags.StringVar(&trackerKey, "tracker-key", trackerKey, "(optional) key used to enable TLS for the tracker")
	flags.StringVar(&trackerCA, "tracker-ca", trackerCA, "(optional) ca used to enable TLS for the tracker")

	flags.StringVar(&tlsKey, "tls-key", tlsKey, "(optional) path to the file containing the TLS private key")
	flags.StringVar(&tlsCert, "tls-cert", tlsCert, "(optional) path to the file containing the TLS certificate")

	err := cmd.Execute()
	exitIff(err)
}
