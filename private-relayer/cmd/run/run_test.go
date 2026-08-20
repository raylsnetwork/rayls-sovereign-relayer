package run

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("returns error when config file not found", func(t *testing.T) {
		err := Run("/nonexistent/path.env")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error reading .env file")
	})

	t.Run("returns error when config validation fails", func(t *testing.T) {
		// Create a temp .env with only one field — all required fields will be missing.
		dir := t.TempDir()
		envFile := filepath.Join(dir, ".env")
		err := os.WriteFile(envFile, []byte("LOG_LEVEL=INFO\n"), 0o600)
		require.NoError(t, err)

		err = Run(envFile)
		require.Error(t, err)
	})
}

func TestLogOnError(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("does nothing when error is nil", func(t *testing.T) {
		logOnError(nil, "should not log anything")
	})

	t.Run("logs when error is not nil", func(t *testing.T) {
		logOnError(errors.New("boom"), "something failed")
	})
}
