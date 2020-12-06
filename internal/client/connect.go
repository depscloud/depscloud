package client

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Connect(cfg *Config) (*grpc.ClientConn, error) {
	options := []grpc.DialOption{
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_prometheus.StreamClientInterceptor,
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_prometheus.UnaryClientInterceptor,
		)),
	}

	if cfg.ServiceConfig != "" {
		options = append(options, grpc.WithDefaultServiceConfig(cfg.ServiceConfig))
	}

	if cfg.TLS || cfg.TLSConfig.CertPath != "" {
		tlsConfig, err := LoadTLSConfig(cfg.TLSConfig)
		if err != nil {
			return nil, err
		}

		tlsCredentials := credentials.NewTLS(tlsConfig)
		options = append(options, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		options = append(options, grpc.WithInsecure())
	}

	return grpc.Dial(cfg.Address, options...)
}
