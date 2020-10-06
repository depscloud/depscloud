package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Connect(cfg *Config) (*grpc.ClientConn, error) {
	options := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(cfg.ServiceConfig),
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
