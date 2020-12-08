package mux_test

import (
	"github.com/depscloud/depscloud/internal/mux"
	"github.com/depscloud/depscloud/internal/v"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig(t *testing.T) {
	in := mux.DefaultConfig(v.Info{})
	out, flags := mux.WithFlags(in)
	require.Equal(t, in, out)
	require.Len(t, flags, 5)
}