// Package mtls builds *tls.Config values for mutual-TLS connections
// between the relayers, CTS, and NATS. Server-side configs require
// and verify a client cert; client-side configs present a cert.
// All peers trust the same CA so any one of them can verify any other.
package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

// LoadServerConfig builds a *tls.Config for a server accepting mTLS
// connections (CTS gRPC). It loads the server keypair, trusts the
// supplied CA for verifying incoming client certs, and requires every
// client to present a cert signed by that CA.
func LoadServerConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("loading server keypair (%s, %s): %w", certFile, keyFile, err))
	}

	pool, err := loadCAPool(caFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// LoadClientConfig builds a *tls.Config for a peer dialling another
// rayls service (CTS gRPC, NATS, etc.). It loads the client keypair,
// trusts the supplied CA for verifying the server certificate, and
// lets the dial target's hostname drive SAN verification.
func LoadClientConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("loading client keypair (%s, %s): %w", certFile, keyFile, err))
	}

	pool, err := loadCAPool(caFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

func loadCAPool(caFile string) (*x509.CertPool, error) {
	pem, err := os.ReadFile(caFile)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("reading CA file %s: %w", caFile, err))
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("CA file %s contains no PEM-encoded certificates", caFile)
	}
	return pool, nil
}
