package mux

import (
	"github.com/depscloud/depscloud/internal/v"

	"github.com/mjpitz/go-gracefully/check"

	"github.com/urfave/cli/v2"

	"google.golang.org/grpc"
)

// ConfigGRPC allows callers to provide custom options to the gRPC server.
type ConfigGRPC struct {
	ServerOptions []grpc.ServerOption
}

// Config defines requirements for starting a server.
type Config struct {
	PortHTTP  int
	PortGRPC  int
	TLSConfig *TLSConfig

	GRPC      *ConfigGRPC
	Checks    []check.Check
	Endpoints []ServerEndpoint
	Version   v.Info
}

// DefaultConfig will construct the default configuration used by projects.
func DefaultConfig(version v.Info) *Config {
	return &Config{
		PortHTTP:  8080,
		PortGRPC:  8090,
		TLSConfig: &TLSConfig{},
		GRPC:      &ConfigGRPC{},
		Version:   version,
	}
}

// WithFlags sets up the appropriate CLI flags for the provided configuration.
func WithFlags(cfg *Config) (*Config, []cli.Flag) {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "http-port",
			Usage:       "the port to run http on",
			Value:       cfg.PortHTTP,
			Destination: &(cfg.PortHTTP),
			EnvVars:     []string{"HTTP_PORT"},
		},
		&cli.IntFlag{
			Name:        "grpc-port",
			Aliases:     []string{"port"},
			Usage:       "the port to run grpc on",
			Value:       cfg.PortGRPC,
			Destination: &(cfg.PortGRPC),
			EnvVars:     []string{"GRPC_PORT"},
		},
		&cli.StringFlag{
			Name:        "tls-key",
			Usage:       "path to the file containing the TLS private key",
			Value:       cfg.TLSConfig.KeyPath,
			Destination: &(cfg.TLSConfig.KeyPath),
			EnvVars:     []string{"TLS_KEY_PATH"},
		},
		&cli.StringFlag{
			Name:        "tls-cert",
			Usage:       "path to the file containing the TLS certificate",
			Value:       cfg.TLSConfig.CertPath,
			Destination: &(cfg.TLSConfig.CertPath),
			EnvVars:     []string{"TLS_CERT_PATH"},
		},
		&cli.StringFlag{
			Name:        "tls-ca",
			Usage:       "path to the file containing the TLS certificate authority",
			Value:       cfg.TLSConfig.CAPath,
			Destination: &(cfg.TLSConfig.CAPath),
			EnvVars:     []string{"TLS_CA_PATH"},
		},
	}

	return cfg, flags
}
