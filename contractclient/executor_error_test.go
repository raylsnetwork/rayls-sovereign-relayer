package contractclient_test

import (
	"strings"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/stretchr/testify/require"
)

func TestErrorWithRevertData_Error(t *testing.T) {
	t.Run("decodes the Error(string) reason so downstream string-matching works", func(t *testing.T) {
		// The Enygma retry classifier matches human-readable revert reasons; the rendered message
		// MUST contain the decoded ASCII reason, not only the raw hex, or those matches never fire.
		err := contractclient.NewErrorWithRevertData(testtools.ErrorStringRevertData(t, "Invalid public signal for balance"))
		msg := err.Error()
		require.Contains(t, msg, "Invalid public signal for balance",
			"Error() must surface the decoded revert reason for downstream string matching")
		// The raw hex is retained for debugging.
		require.Contains(t, msg, "0x08c379a0")
	})

	t.Run("falls back to hex for custom-error reverts that don't decode to a string", func(t *testing.T) {
		// ProgramData__UnapprovedTemplate(uint256,bytes32,bytes4) selector 0xccdb0e9f — not Error(string).
		custom := []byte{0xcc, 0xdb, 0x0e, 0x9f, 0x00, 0x01, 0x02, 0x03}
		msg := contractclient.NewErrorWithRevertData(custom).Error()
		require.True(t, strings.HasPrefix(msg, "transaction reverted: 0x"),
			"non-Error(string) reverts keep the hex payload; got %q", msg)
		require.Contains(t, msg, "ccdb0e9f")
	})

	t.Run("empty revert data", func(t *testing.T) {
		require.Equal(t, "transaction reverted", contractclient.NewErrorWithRevertData(nil).Error())
	})
}
