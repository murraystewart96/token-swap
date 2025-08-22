package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/events"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createEventListenerCmd() *cobra.Command { //nolint:gocognit
	clientCmd := &cobra.Command{
		Use:   "event-listener",
		Short: "event listener",
		Long:  `Listens for transaction events and published to event queue`,

		Run: func(cmd *cobra.Command, _ []string) {
			// get command flags
			configPath, err := cmd.Flags().GetString(configFlag)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse command flag")
			}

			cfg := &config.Events{}
			config.ReadEnvironment(configPath, cfg)

			events, err := events.NewClient(cfg)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create event service")
			}

			ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

			events.Listen(ctx)
		},
	}

	return clientCmd
}
