package client_test

import (
	"testing"

	"github.com/depscloud/depscloud/internal/client"

	"github.com/stretchr/testify/require"
)

func Test_WithFlags(t *testing.T) {
	in := &client.Config{
		TLSConfig: &client.TLSConfig{},
	}
	out, flags := client.WithFlags("extractor", in)
	require.Equal(t, in, out)
	require.Len(t, flags, 7)
}
