package mux

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

func LoadTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	if cfg == nil || cfg.CertPath == "" || cfg.KeyPath == "" {
		return nil, nil
	}

	certificate, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
	if err != nil {
		return nil, err
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

	return &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}, nil
}
