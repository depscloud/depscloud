package logger

import (
	"github.com/urfave/cli/v2"

	"go.uber.org/zap/zapcore"
)

// wrapper for setting zap's log level through urfave CLI

type logLevelWrapper struct {
	cfg *Config
}

func (l *logLevelWrapper) Set(value string) error {
	var level zapcore.Level
	if err := level.Set(value); err != nil {
		return err
	}
	l.cfg.ZapConfig.Level.SetLevel(level)
	return nil
}

func (l *logLevelWrapper) String() string {
	return l.cfg.ZapConfig.Level.Level().String()
}

var _ cli.Generic = &logLevelWrapper{}
