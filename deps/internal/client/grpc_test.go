package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_translateBaseURL(t *testing.T) {
	{
		tls, translated := translateBaseURL(DefaultBaseURL)
		require.True(t, tls)
		require.Equal(t, "api.deps.cloud:443", translated)
	}

	{
		tls, translated := translateBaseURL("https://api.deps.cloud:1234")
		require.True(t, tls)
		require.Equal(t, "api.deps.cloud:1234", translated)
	}

	{
		tls, translated := translateBaseURL("client://api.deps.cloud")
		require.False(t, tls)
		require.Equal(t, "api.deps.cloud:80", translated)
	}

	{
		tls, translated := translateBaseURL("client://api.deps.cloud:1234")
		require.False(t, tls)
		require.Equal(t, "api.deps.cloud:1234", translated)
	}
}
