package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/events"
	"github.com/murraystewart96/token-swap/pkg/kafka"
	"github.com/murraystewart96/token-swap/internal/storage/postgres"
	"github.com/murraystewart96/token-swap/pkg/tracing"
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

			// Initialize tracing
			tracingConfig := tracing.TracingConfig{
				ServiceName:  "token-swap-event-listener",
				Environment:  "development", // TODO: Make this configurable
				OTLPEndpoint: "", // Will use default for development
			}
			shutdown, err := tracing.InitTracer(tracingConfig)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to initialize tracing")
			}
			defer shutdown()

			producer, err := kafka.NewProducer(&cfg.Kafka)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create kafka producer")
			}

			db, err := postgres.NewDB(&cfg.DB)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create database connection")
			}
			defer db.Close()

			events, err := events.NewClient(&cfg.Listener, producer, db)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create event service")
			}

			ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

			events.Listen(ctx)
		},
	}

	return eventsCmd
}
