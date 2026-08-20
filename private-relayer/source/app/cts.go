package app

import (
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/txops"
)

type CTSClient struct {
	keys.KeysServiceClient
	encrypt.EncryptServiceClient

	txops.PrivateHubTxOpsServiceClient
	txops.PrivateNodeTxOpsServiceClient
	txops.DVPOperatorTxOpsServiceClient
}

func (r *SourcePrivateRelayer) initializeCTSClients() {
	r.ctsClient = &CTSClient{
		KeysServiceClient:    keys.NewKeysServiceClient(r.ctsConn),
		EncryptServiceClient: encrypt.NewEncryptServiceClient(r.ctsConn),

		PrivateHubTxOpsServiceClient:  txops.NewPrivateHubTxOpsServiceClient(r.ctsConn),
		PrivateNodeTxOpsServiceClient: txops.NewPrivateNodeTxOpsServiceClient(r.ctsConn),
		DVPOperatorTxOpsServiceClient: txops.NewDVPOperatorTxOpsServiceClient(r.ctsConn),
	}
}
