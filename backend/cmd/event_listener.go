package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/events"
	"github.com/murraystewart96/token-swap/internal/kafka"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createEventListenerCmd() *cobra.Command { //nolint:gocognit
	eventsCmd := &cobra.Command{
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

			producer, err := kafka.NewProducer(&cfg.Kafka)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create kafka producer")
			}

			events, err := events.NewClient(&cfg.Listener, producer)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create event service")
			}

			ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

			events.Listen(ctx)
		},
	}

	return eventsCmd
}
