package logger_test

import (
	"testing"

	"github.com/depscloud/depscloud/internal/logger"

	"github.com/stretchr/testify/require"
)

func Test_WithFlags(t *testing.T) {
	in := logger.DefaultConfig()
	out, flags := logger.WithFlags(in)
	require.Equal(t, in, out)
	require.Len(t, flags, 2)
}
