package mux_test

import (
	"testing"

	"github.com/depscloud/depscloud/internal/appconf"
	"github.com/depscloud/depscloud/internal/mux"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	in := mux.DefaultConfig(&appconf.V{})
	out, flags := mux.WithFlags(in)
	require.Equal(t, in, out)
	require.Len(t, flags, 5)
}
