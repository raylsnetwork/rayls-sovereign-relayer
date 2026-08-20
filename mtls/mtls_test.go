package mtls_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/mtls"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/health"
)

// TestMTLSHandshake validates that the configs produced by LoadServerConfig
// and LoadClientConfig complete a real gRPC mTLS handshake when fed certs
// signed by the same CA — and that a client cert signed by a foreign CA is
// rejected by the server.
func TestMTLSHandshake(t *testing.T) {
	dir := t.TempDir()

	caCert, caKey := newCA(t)
	writePEM(t, filepath.Join(dir, "ca.crt"), "CERTIFICATE", caCert.Raw)

	serverCert, serverKey := newLeaf(t, caCert, caKey, "127.0.0.1", []string{"localhost"}, []net.IP{net.ParseIP("127.0.0.1")}, true)
	writePEM(t, filepath.Join(dir, "server.crt"), "CERTIFICATE", serverCert.Raw)
	writeKey(t, filepath.Join(dir, "server.key"), serverKey)

	clientCert, clientKey := newLeaf(t, caCert, caKey, "client", nil, nil, false)
	writePEM(t, filepath.Join(dir, "client.crt"), "CERTIFICATE", clientCert.Raw)
	writeKey(t, filepath.Join(dir, "client.key"), clientKey)

	serverTLS, err := mtls.LoadServerConfig(
		filepath.Join(dir, "ca.crt"),
		filepath.Join(dir, "server.crt"),
		filepath.Join(dir, "server.key"),
	)
	require.NoError(t, err)

	clientTLS, err := mtls.LoadClientConfig(
		filepath.Join(dir, "ca.crt"),
		filepath.Join(dir, "client.crt"),
		filepath.Join(dir, "client.key"),
	)
	require.NoError(t, err)

	addr := startHealthServer(t, serverTLS)

	t.Run("trusted client succeeds", func(t *testing.T) {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(clientTLS)))
		require.NoError(t, err)
		t.Cleanup(func() { _ = conn.Close() })

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = healthgrpc.NewHealthClient(conn).Check(ctx, &healthgrpc.HealthCheckRequest{})
		require.NoError(t, err)
	})

	t.Run("foreign-CA client is rejected", func(t *testing.T) {
		foreignCA, foreignKey := newCA(t)
		intruderCert, intruderKey := newLeaf(t, foreignCA, foreignKey, "intruder", nil, nil, false)
		fdir := t.TempDir()
		writePEM(t, filepath.Join(fdir, "ca.crt"), "CERTIFICATE", foreignCA.Raw)
		writePEM(t, filepath.Join(fdir, "client.crt"), "CERTIFICATE", intruderCert.Raw)
		writeKey(t, filepath.Join(fdir, "client.key"), intruderKey)

		intruderTLS, err := mtls.LoadClientConfig(
			filepath.Join(dir, "ca.crt"), // trust real CA so server cert verifies
			filepath.Join(fdir, "client.crt"),
			filepath.Join(fdir, "client.key"),
		)
		require.NoError(t, err)

		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(intruderTLS)))
		require.NoError(t, err)
		t.Cleanup(func() { _ = conn.Close() })

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err = healthgrpc.NewHealthClient(conn).Check(ctx, &healthgrpc.HealthCheckRequest{})
		require.Error(t, err, "server must reject client signed by foreign CA")
	})
}

func TestLoadServerConfig_MissingFile(t *testing.T) {
	_, err := mtls.LoadServerConfig("/no/such/ca", "/no/such/cert", "/no/such/key")
	require.Error(t, err)
}

func TestLoadClientConfig_BadCAPEM(t *testing.T) {
	dir := t.TempDir()

	caCert, caKey := newCA(t)
	cert, key := newLeaf(t, caCert, caKey, "x", nil, nil, false)
	writePEM(t, filepath.Join(dir, "x.crt"), "CERTIFICATE", cert.Raw)
	writeKey(t, filepath.Join(dir, "x.key"), key)

	badCA := filepath.Join(dir, "bad.crt")
	require.NoError(t, os.WriteFile(badCA, []byte("not a pem"), 0o644))

	_, err := mtls.LoadClientConfig(badCA, filepath.Join(dir, "x.crt"), filepath.Join(dir, "x.key"))
	require.Error(t, err)
}

// --- helpers ---

func startHealthServer(t *testing.T, conf *tls.Config) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv := grpc.NewServer(grpc.Creds(credentials.NewTLS(conf)))
	healthgrpc.RegisterHealthServer(srv, health.NewServer())

	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	return lis.Addr().String()
}

func newCA(t *testing.T) (*x509.Certificate, *rsa.PrivateKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)
	cert, err := x509.ParseCertificate(der)
	require.NoError(t, err)
	return cert, key
}

func newLeaf(t *testing.T, ca *x509.Certificate, caKey *rsa.PrivateKey, cn string, dnsNames []string, ips []net.IP, server bool) (*x509.Certificate, *rsa.PrivateKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	usage := x509.ExtKeyUsageClientAuth
	if server {
		usage = x509.ExtKeyUsageServerAuth
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{usage},
		DNSNames:     dnsNames,
		IPAddresses:  ips,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca, &key.PublicKey, caKey)
	require.NoError(t, err)
	cert, err := x509.ParseCertificate(der)
	require.NoError(t, err)
	return cert, key
}

func writePEM(t *testing.T, path, kind string, der []byte) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: kind, Bytes: der}), 0o644))
}

func writeKey(t *testing.T, path string, key *rsa.PrivateKey) {
	t.Helper()
	der := x509.MarshalPKCS1PrivateKey(key)
	require.NoError(t, os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}), 0o600))
}
