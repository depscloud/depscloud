package logger

import (
	"github.com/urfave/cli/v2"

	"go.uber.org/zap"
)

type Config struct {
	ZapConfig zap.Config
}

// DefaultConfig returns the default configuration used to
func DefaultConfig() *Config {
	return &Config{
		ZapConfig: zap.NewProductionConfig(),
	}
}

// WithFlags configures flags for the provided configuration.
func WithFlags(cfg *Config) (*Config, []cli.Flag) {
	logLevel := &logLevelWrapper{
		cfg: cfg,
	}

	flags := []cli.Flag{
		&cli.GenericFlag{
			Name:    "log-level",
			Usage:   "configures the level at with logs are written",
			Value:   logLevel,
			EnvVars: []string{"LOG_LEVEL"},
		},
		&cli.StringFlag{
			Name:        "log-format",
			Usage:       "configures the format of the logs (console / json)",
			Value:       cfg.ZapConfig.Encoding,
			Destination: &(cfg.ZapConfig.Encoding),
			EnvVars:     []string{"LOG_FORMAT"},
		},
	}

	return cfg, flags
}

// MustGetLogger uses the provided configuration to construct a logger.
// If an error occurs, it panics.
func MustGetLogger(cfg *Config) *zap.Logger {
	logger, err := cfg.ZapConfig.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
