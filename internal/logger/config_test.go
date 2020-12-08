package logger_test

import (
	"github.com/depscloud/depscloud/internal/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func Test_WithFlags(t *testing.T) {
	in := zap.NewProductionConfig()
	_, flags := logger.WithFlags(in)
	require.Len(t, flags, 2)
}