package grpc

import (
	"fmt"
	"net"

	encryptpb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/encrypt"
	keyspb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	txopspb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/txops"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/mtls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TLSConfig points at the PEM-encoded files used to authenticate this
// server and verify incoming client certs. All three are required —
// the CTS gRPC channel runs mTLS unconditionally.
type TLSConfig struct {
	CAFile   string
	CertFile string
	KeyFile  string
}

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

// NewServer wires the CTS gRPC server. privateChainTxOpsHandler and
// publicChainTxOpsHandler are optional: when nil the corresponding TxOps
// service is not registered and any RPC against it receives gRPC
// Unimplemented. Both are skipped together when CTS runs without the public
// chain stack (no pubrelayer-* services to call them). All other handlers are
// required.
func NewServer(
	addr string,
	tlsCfg TLSConfig,
	encryptHandler *EncryptHandler,
	keysHandler *KeysHandler,

	privateHubTxOpsHandler *TxOpsHandler,
	privateNodeTxOpsHandler *TxOpsHandler,
	dvpOperatorTxOpsHandler *TxOpsHandler,

	privateChainTxOpsHandler *TxOpsHandler,
	publicChainTxOpsHandler *TxOpsHandler,
) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", addr, err)
	}

	tlsConf, err := mtls.LoadServerConfig(tlsCfg.CAFile, tlsCfg.CertFile, tlsCfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("loading server TLS config: %w", err)
	}

	srv := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConf)))
	encryptpb.RegisterEncryptServiceServer(srv, encryptHandler)
	keyspb.RegisterKeysServiceServer(srv, keysHandler)

	txopspb.RegisterPrivateNodeTxOpsServiceServer(srv, privateNodeTxOpsHandler)

	if privateHubTxOpsHandler != nil {
		txopspb.RegisterPrivateHubTxOpsServiceServer(srv, privateHubTxOpsHandler)
	}
	if dvpOperatorTxOpsHandler != nil {
		txopspb.RegisterDVPOperatorTxOpsServiceServer(srv, dvpOperatorTxOpsHandler)
	}

	if privateChainTxOpsHandler != nil {
		txopspb.RegisterPrivateChainTxOpsServiceServer(srv, privateChainTxOpsHandler)
	}
	if publicChainTxOpsHandler != nil {
		txopspb.RegisterPublicChainTxOpsServiceServer(srv, publicChainTxOpsHandler)
	}

	return &Server{
		grpcServer: srv,
		listener:   lis,
	}, nil
}

func (s *Server) Serve() error {
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}
