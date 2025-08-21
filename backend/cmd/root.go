package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	// flags.
	configFlag = "config"
)

func createRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "token-swap",
		Short: `Backend services for token-swap dApp`,
		Long:  `Backend services for token-swap dApp including event listener, worker, and API server`,
	}

	var configPath string

	rootCmd.PersistentFlags().StringVar(&configPath, configFlag, "", "Path to config including file name")

	return rootCmd
}

func addSubcmds(rootCmd *cobra.Command) {
	rootCmd.AddCommand(createEventListenerCmd())
}

func Execute() {
	rootCmd := createRootCmd()
	addSubcmds(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute root command")
	}
}
