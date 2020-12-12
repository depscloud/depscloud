package proxy_test

import (
	"fmt"
	"testing"

	"github.com/depscloud/api/v1alpha/extractor"
	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/depscloud/gateway/internal/proxy"

	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
)

func first(in map[string]grpc.ServiceInfo) (string, grpc.ServiceInfo) {
	for key, val := range in {
		return key, val
	}
	return "", grpc.ServiceInfo{}
}

func TestRouter(t *testing.T) {
	aServer := &grpc.ClientConn{}
	var aSnapshot map[string]grpc.ServiceInfo

	bServer := &grpc.ClientConn{}
	var bSnapshot map[string]grpc.ServiceInfo

	router, err := proxy.NewRouter([]*proxy.Backend{
		{
			ClientConn: aServer,
			RegisterService: func(server *grpc.Server) {
				extractor.RegisterDependencyExtractorServer(server, &extractor.UnimplementedDependencyExtractorServer{})

				aSnapshot = server.GetServiceInfo()
			},
		},
		{
			ClientConn: bServer,
			RegisterService: func(server *grpc.Server) {
				extractor.RegisterDependencyExtractorServer(server, &extractor.UnimplementedDependencyExtractorServer{})
				v1beta.RegisterManifestExtractionServiceServer(server, &v1beta.UnimplementedManifestExtractionServiceServer{})

				bSnapshot = server.GetServiceInfo()
			},
		},
	}...)
	require.NoError(t, err)

	// backends are run
	require.Len(t, aSnapshot, 1)
	require.Len(t, bSnapshot, 2)

	// verify v1alpha is routed to aServer, despite the service registration in bServer.
	{
		serviceName, _ := first(aSnapshot)
		fullMethodName := fmt.Sprintf("/%s/TestMethod", serviceName)

		cc, err := router.Route(fullMethodName)
		require.NoError(t, err)
		require.Equal(t, aServer, cc)

		// remove from the other snapshot
		delete(bSnapshot, serviceName)
	}

	// verify v1beta is routed to bServer
	{
		serviceName, _ := first(bSnapshot)
		fullMethodName := fmt.Sprintf("/%s/TestMethod", serviceName)

		cc, err := router.Route(fullMethodName)
		require.NoError(t, err)
		require.Equal(t, bServer, cc)
	}
}
