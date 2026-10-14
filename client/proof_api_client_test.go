package client_test

import (
	"errors"
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-relayer/client"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type failingHTTPClient struct{ calls int }

func (f *failingHTTPClient) Do(*http.Request) (*http.Response, error) {
	f.calls++
	return nil, errors.New("unexpected HTTP call")
}

func TestProofAPIClient_CreateWithdrawProof_MissingPaymentOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(req *types.WithdrawProofRequest)
	}{
		{name: "nil payment commitment", mutate: func(req *types.WithdrawProofRequest) { req.PaymentCommitment = nil }},
		{name: "nil payment secret key", mutate: func(req *types.WithdrawProofRequest) { req.PaymentSecretKey = nil }},
		{name: "nil payment salt", mutate: func(req *types.WithdrawProofRequest) { req.PaymentSalt = nil }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			httpClient := &failingHTTPClient{}
			proofClient := client.NewProofAPIClient("http://proof.invalid", httpClient)
			req := types.WithdrawProofRequest{
				PaymentCommitment: big.NewInt(1),
				PaymentSecretKey:  big.NewInt(2),
				PaymentSalt:       big.NewInt(3),
			}
			tc.mutate(&req)

			require.NotPanics(t, func() {
				_, err := proofClient.CreateWithdrawProof(2, req)
				assert.ErrorContains(t, err, "withdraw proof request is missing a payment-output scalar")
			})
			assert.Zero(t, httpClient.calls)
		})
	}
}
