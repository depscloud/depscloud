package client

import (
	"strings"

	"github.com/urfave/cli/v2"
)

// https://github.com/grpc/grpc/blob/master/doc/service_config.md
const DefaultServiceConfig = `{"loadBalancingPolicy":"round_robin","healthCheckConfig":{"serviceName":""}}`

const DefaultLoadBalancer = "round_robin"

type Config struct {
	Address       string
	ServiceConfig string
	LoadBalancer  string
	TLS           bool
	TLSConfig     *TLSConfig
}

func WithFlags(prefix string, cfg *Config) (*Config, []cli.Flag) {
	lower := strings.ToLower(prefix)
	upper := strings.ToUpper(prefix)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        lower + "-address",
			Usage:       "address to the " + lower,
			Value:       cfg.Address,
			Destination: &(cfg.Address),
			EnvVars:     []string{upper + "_ADDRESS"},
		},
		&cli.StringFlag{
			Name:        lower + "-service-config",
			Usage:       "service configuration for the " + lower,
			Value:       cfg.ServiceConfig,
			Destination: &(cfg.ServiceConfig),
			EnvVars:     []string{upper + "_SERVICE_CONFIG"},
		},
		&cli.BoolFlag{
			Name:        lower + "-tls",
			Usage:       "enable TLS for the " + lower,
			Value:       cfg.TLS,
			Destination: &(cfg.TLS),
			EnvVars:     []string{upper + "_TLS"},
		},
		&cli.StringFlag{
			Name:        lower + "-ca",
			Usage:       "ca used to enable TLS for the " + lower,
			Value:       cfg.TLSConfig.CAPath,
			Destination: &(cfg.TLSConfig.CAPath),
			EnvVars:     []string{upper + "_CA_PATH"},
		},
		&cli.StringFlag{
			Name:        lower + "-cert",
			Usage:       "certificate used to enable TLS for the " + lower,
			Value:       cfg.TLSConfig.CertPath,
			Destination: &(cfg.TLSConfig.CertPath),
			EnvVars:     []string{upper + "_CERT_PATH"},
		},
		&cli.StringFlag{
			Name:        lower + "-key",
			Usage:       "key used to enable TLS for the " + lower,
			Value:       cfg.TLSConfig.KeyPath,
			Destination: &(cfg.TLSConfig.KeyPath),
			EnvVars:     []string{upper + "_KEY_PATH"},
		},
		// deprecated
		&cli.StringFlag{
			Name:        lower + "-lb",
			Usage:       "the load balancer policy to use for the " + lower,
			Value:       cfg.LoadBalancer,
			Destination: &(cfg.LoadBalancer),
			EnvVars:     []string{upper + "_LBPOLICY"},
		},
	}

	return cfg, flags
}
