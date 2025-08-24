package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/kafka"
	"github.com/murraystewart96/token-swap/internal/worker"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createWorkerCmd() *cobra.Command { //nolint:gocognit
	workerCmd := &cobra.Command{
		Use:   "worker",
		Short: "processes events",
		Long:  `Consumes events produced by the event-listener`,

		Run: func(cmd *cobra.Command, _ []string) {
			// get command flags
			configPath, err := cmd.Flags().GetString(configFlag)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse command flag")
			}

			cfg := &config.Worker{}
			config.ReadEnvironment(configPath, cfg)

			tradeHistoryConsumer, err := kafka.NewConsumer(&cfg.Kafka)

			worker := worker.New(tradeHistoryConsumer)

			ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
			worker.Start(ctx)
		},
	}

	return workerCmd
}
