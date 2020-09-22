package mux

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

type TLSConfig struct {
	CertPath string
	KeyPath string
	CAPath string
}

func LoadTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	if cfg.CertPath == "" || cfg.KeyPath == "" {
		return nil, nil
	}

	certificate, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if cfg.CAPath != "" {
		bs, err := ioutil.ReadFile(cfg.CAPath)
		if err != nil {
			return nil, err
		}

		ok := certPool.AppendCertsFromPEM(bs)
		if !ok {
			return nil, fmt.Errorf("failed to append certs")
		}
	}

	return &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}, nil
}
