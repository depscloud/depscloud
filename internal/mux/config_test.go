package mux_test

import (
	"testing"

	"github.com/depscloud/depscloud/internal/mux"
	"github.com/depscloud/depscloud/internal/v"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	in := mux.DefaultConfig(v.Info{})
	out, flags := mux.WithFlags(in)
	require.Equal(t, in, out)
	require.Len(t, flags, 5)
}
