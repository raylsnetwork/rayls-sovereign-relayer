// Decommissioning Teleport (vanilla, atomic).

// Package main implements the legacy public-chain (RN) Teleport bridge relayer.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/cmd/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags for running pubrelayer app
var EnvPathFlagName = "env"

var rootCmd = &cobra.Command{
	Use:   "rayl-public-relayer",
	Short: "A command-line tool for rayls-public-relayer",
	Long:  `rayls-public-relayer is a command-line tool for deploying and running rayls-public-relayer.`,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run public relayer app",
	Long:  "Run public relayer app",
	RunE: func(cmd *cobra.Command, args []string) error {
		envPath := viper.GetString(EnvPathFlagName)

		if err := run.Run(envPath); err != nil {
			slog.Error("Error running the relayer", slog.Any("error", err))
			return err
		}
		return nil
	},
}

func bindFlags(rootCMD *cobra.Command) {
	rootCMD.PersistentFlags().String(EnvPathFlagName, "", "Path to .env file")
	_ = viper.BindPFlag(EnvPathFlagName, rootCMD.PersistentFlags().Lookup(EnvPathFlagName))
}

func init() {
	bindFlags(runCmd)
	rootCmd.AddCommand(runCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
