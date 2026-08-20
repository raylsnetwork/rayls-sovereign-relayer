package main

import (
	"fmt"
	"os"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/cmd/migrate"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/cmd/run"
	"github.com/spf13/cobra"
)

var (
	envPath         string
	EnvPathFlagName = "env"

	fromEncryptor         string
	FromEncryptorFlagName = "from"
	toEncryptor           string
	ToEncryptorFlagName   = "to"
)

var rootCmd = &cobra.Command{
	Use:   "rayls-cts",
	Short: "A command-line tool for Rayls CTS service",
	Long:  `This is a command-line tool for running the Rayls Cryptography Trust Suite`,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Rayls CTS service",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		if err := run.RunServer(envPath); err != nil {
			return fmt.Errorf("running CTS server: %w", err)
		}
		return nil
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate keys from one encryption service to another",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return migrate.Migrate(envPath, fromEncryptor, toEncryptor)
	},
}

func BindFlags(rootCMD *cobra.Command) {
	rootCMD.PersistentFlags().StringVar(&envPath, EnvPathFlagName, "", "Path to .env file")

	migrateCmd.Flags().
		StringVar(&fromEncryptor, FromEncryptorFlagName, "", "Encryption service to migrate the database from")
	_ = migrateCmd.MarkFlagRequired(FromEncryptorFlagName)
	migrateCmd.Flags().StringVar(&toEncryptor, ToEncryptorFlagName, "", "Encryption service to migrate the database to")
	_ = migrateCmd.MarkFlagRequired(ToEncryptorFlagName)
}

func main() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(migrateCmd)
	BindFlags(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
