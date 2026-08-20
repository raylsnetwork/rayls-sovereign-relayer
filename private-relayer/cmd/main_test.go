package main

import (
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCobraCommandStructure(t *testing.T) {
	t.Run("run command is registered as subcommand", func(t *testing.T) {
		var found bool
		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == "run" {
				found = true
				break
			}
		}
		assert.True(t, found, "run subcommand should be registered on rootCmd")
	})

	t.Run("run command has env flag", func(t *testing.T) {
		flag := runCmd.PersistentFlags().Lookup(EnvPathFlagName)
		require.NotNil(t, flag, "--env flag should be registered")
		assert.Equal(t, "", flag.DefValue)
	})

	t.Run("run command returns error for invalid env path", func(t *testing.T) {
		testtools.SilenceLogger()

		rootCmd.SetArgs([]string{"run", "--env", "/nonexistent/path.env"})
		err := rootCmd.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error reading .env file")
	})
}
