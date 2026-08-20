package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/cmd/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags for running the Raylz Relayer app
var EnvPathFlagName = "env"

var rootCmd = &cobra.Command{
	Use:   "rayls-relayer",
	Short: "A command-line tool for rayls-relayer",
	Long:  `rayls-relayer is a command-line tool for deploying and running rayls-relayer.`,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run rayls-relayer",
	RunE: func(cmd *cobra.Command, args []string) error {
		envPath := viper.GetString(EnvPathFlagName)
		if err := run.Run(envPath); err != nil {
			slog.Error("Error running the relayer", slog.Any("error", err))
			return fmt.Errorf("running private relayer: %w", err)
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
