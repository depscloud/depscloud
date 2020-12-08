package logger_test

import (
	"context"
	"github.com/depscloud/depscloud/internal/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestContext(t *testing.T) {
	expected, err := zap.NewProduction()
	require.NoError(t, err)

	ctx := logger.ToContext(context.Background(), expected)
	actual := logger.Extract(ctx)

	require.Equal(t, expected, actual)
}
