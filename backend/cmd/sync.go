package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/storage/postgres"
	"github.com/murraystewart96/token-swap/internal/storage/redis"
	"github.com/murraystewart96/token-swap/internal/sync"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createSyncCmd() *cobra.Command { //nolint:gocognit
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "syncs database with contract",
		Long:  `Sync database and cache with contract's state on the blockchain`,

		Run: func(cmd *cobra.Command, _ []string) {
			// get command flags
			configPath, err := cmd.Flags().GetString(configFlag)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse command flag")
			}

			cfg := &config.Sync{}
			config.ReadEnvironment(configPath, cfg)

			poolCache := redis.NewCache(&cfg.Redis)
			db, err := postgres.NewDB(&cfg.DB)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to connect to DB")
			}

			syncClient, err := sync.NewSync(cfg, poolCache, db)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create Sync client")
			}

			ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
			syncClient.Start(ctx)
		},
	}

	return syncCmd
}
