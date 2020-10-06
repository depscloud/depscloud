package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

type TLSConfig struct {
	CertPath string
	KeyPath  string
	CAPath   string
}

func LoadTLSConfig(cfg *TLSConfig) (tlsConfig *tls.Config, err error) {
	if cfg == nil {
		return nil, nil
	}

	tlsConfig = &tls.Config{}

	if cfg.CertPath != "" && cfg.KeyPath != "" {
		certificate, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{certificate}
	}

	var certPool *x509.CertPool

	if cfg.CAPath != "" {
		caPEM, err := ioutil.ReadFile(cfg.CAPath)
		if err != nil {
			return nil, err
		}

		certPool = x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM(caPEM)
		if !ok {
			return nil, fmt.Errorf("failed to append certs")
		}
	} else {
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
	}

	tlsConfig.ClientCAs = certPool

	return tlsConfig, nil
}
