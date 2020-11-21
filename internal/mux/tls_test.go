package mux_test

import (
	"crypto/tls"
	"github.com/depscloud/depscloud/internal/mux"
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

func Test_MuxLoadTLSConfig(t *testing.T) {
	{
		cfg, err := mux.LoadTLSConfig(nil)
		require.NoError(t, err)
		require.Nil(t, cfg)
	}

	{
		cfg, err := mux.LoadTLSConfig(&mux.TLSConfig{
			CertPath: "missing_key.crt",
		})
		require.NoError(t, err)
		require.Nil(t, cfg)
	}

	{
		cfg, err := mux.LoadTLSConfig(&mux.TLSConfig{
			KeyPath: "missing_cert.key",
		})
		require.NoError(t, err)
		require.Nil(t, cfg)
	}

	{
		cfg, err := mux.LoadTLSConfig(&mux.TLSConfig{
			CertPath: "nonexistent_cert.crt",
			KeyPath: "nonexistent_key.key",
		})
		require.Error(t, err)
		require.Nil(t, cfg)
	}

	{
		cfg, err := mux.LoadTLSConfig(&mux.TLSConfig{
			CertPath: path.Join("..", "hack", "test.crt"),
			KeyPath: path.Join("..", "hack", "test.key"),
		})
		require.NoError(t, err)
		require.NotNil(t, cfg)

		require.Equal(t, cfg.ClientAuth, tls.RequireAndVerifyClientCert)
		require.Len(t, cfg.Certificates, 1)
	}

	{
		cfg, err := mux.LoadTLSConfig(&mux.TLSConfig{
			CAPath: path.Join("..", "hack", "ca.crt"),
			CertPath: path.Join("..", "hack", "test.crt"),
			KeyPath: path.Join("..", "hack", "test.key"),
		})
		require.NoError(t, err)
		require.NotNil(t, cfg)

		require.Equal(t, cfg.ClientAuth, tls.RequireAndVerifyClientCert)
		require.Len(t, cfg.Certificates, 1)
	}
}
