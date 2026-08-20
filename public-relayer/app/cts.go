// Decommissioning Teleport (vanilla, atomic).

package app

import (
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/txops"
)

type CTSClients struct {
	txops.PrivateChainTxOpsServiceClient
	txops.PublicChainTxOpsServiceClient
}

func (r *PublicRelayer) initializeCTSClients() {
	r.ctsClients = &CTSClients{
		PrivateChainTxOpsServiceClient: txops.NewPrivateChainTxOpsServiceClient(r.ctsConn),
		PublicChainTxOpsServiceClient:  txops.NewPublicChainTxOpsServiceClient(r.ctsConn),
	}
}
