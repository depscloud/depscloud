package client_test

import (
	"github.com/depscloud/depscloud/internal/client"
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

func Test_ClientLoadTLSConfig(t *testing.T) {
	{
		cfg, err := client.LoadTLSConfig(nil)
		require.NoError(t, err)
		require.Nil(t, cfg)
	}

	{
		cfg, err := client.LoadTLSConfig(&client.TLSConfig{})
		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.Len(t, cfg.Certificates, 0)
	}

	{
		cfg, err := client.LoadTLSConfig(&client.TLSConfig{
			CertPath: "missing_key.crt",
		})
		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.Len(t, cfg.Certificates, 0)
	}

	{
		cfg, err := client.LoadTLSConfig(&client.TLSConfig{
			KeyPath: "missing_cert.key",
		})
		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.Len(t, cfg.Certificates, 0)
	}

	{
		cfg, err := client.LoadTLSConfig(&client.TLSConfig{
			CertPath: "nonexistent_cert.crt",
			KeyPath: "nonexistent_key.key",
		})
		require.Error(t, err)
		require.Nil(t, cfg)
	}

	{
		cfg, err := client.LoadTLSConfig(&client.TLSConfig{
			CertPath: path.Join("..", "hack", "test.crt"),
			KeyPath: path.Join("..", "hack", "test.key"),
		})
		require.NoError(t, err)
		require.NotNil(t, cfg)

		require.Len(t, cfg.Certificates, 1)
	}

	{
		cfg, err := client.LoadTLSConfig(&client.TLSConfig{
			CAPath: path.Join("..", "hack", "ca.crt"),
			CertPath: path.Join("..", "hack", "test.crt"),
			KeyPath: path.Join("..", "hack", "test.key"),
		})
		require.NoError(t, err)
		require.NotNil(t, cfg)

		require.Len(t, cfg.Certificates, 1)
	}
}
