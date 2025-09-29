package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/kafka"
	"github.com/murraystewart96/token-swap/internal/storage/postgres"
	"github.com/murraystewart96/token-swap/internal/storage/redis"
	"github.com/murraystewart96/token-swap/internal/worker"
	"github.com/murraystewart96/token-swap/pkg/tracing"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createWorkerCmd() *cobra.Command { //nolint:gocognit
	workerCmd := &cobra.Command{
		Use:   "worker",
		Short: "processes events",
		Long:  `Consumes events produced by the event-listener`,

		Run: func(cmd *cobra.Command, _ []string) {
			// Read config
			configPath, err := cmd.Flags().GetString(configFlag)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse command flag")
			}
			cfg := &config.Worker{}
			config.ReadEnvironment(configPath, cfg)

			// Initialize tracing
			tracingConfig := tracing.TracingConfig{
				ServiceName:  "token-swap-worker",
				Environment:  "development", // TODO: Make this configurable
				OTLPEndpoint: "", // Will use default for development
			}
			shutdown, err := tracing.InitTracer(tracingConfig)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to initialize tracing")
			}
			defer shutdown()

			// Worker dependencies
			consumer, err := kafka.NewConsumer(&cfg.Kafka)
			poolCache := redis.NewCache(&cfg.Redis)
			db, err := postgres.NewDB(&cfg.DB)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to connect to DB")
			}

			// Create and start worker
			worker, err := worker.New(consumer, cfg.Topics, poolCache, db)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create worker")
			}

			ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
			worker.Start(ctx)
		},
	}

	return workerCmd
}
