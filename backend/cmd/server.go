package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/server"
	"github.com/murraystewart96/token-swap/internal/storage/postgres"
	"github.com/murraystewart96/token-swap/internal/storage/redis"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createServerCmd() *cobra.Command { //nolint:gocognit
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "api server",
		Long:  `Starts api server`,

		Run: func(cmd *cobra.Command, _ []string) {
			// get command flags
			configPath, err := cmd.Flags().GetString(configFlag)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse command flag")
			}

			cfg := &config.Server{}
			config.ReadEnvironment(configPath, cfg)

			poolCache := redis.NewCache(&cfg.Redis)
			db, err := postgres.NewDB(&cfg.DB)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to connect to DB")
			}

			ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

			handler := server.NewHandler(db, poolCache)
			server.Serve(ctx, cfg.Addr, handler)
		},
	}

	return serverCmd
}
