package cmd

import (
	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/storage/postgres"
	"github.com/murraystewart96/token-swap/internal/storage/postgres/migrations"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createMigrateCmd() *cobra.Command {
	migrationsCmd := &cobra.Command{
		Use:   "migrate",
		Short: `Runs the db migrations`,
		PreRun: func(_ *cobra.Command, _ []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
			// get command flags
			configPath, err := cmd.Flags().GetString(configFlag)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse command flag")
			}

			cfg := &config.Migration{}
			config.ReadEnvironment(configPath, cfg)

			log.Debug().Interface("config", cfg).Msg("Parsed config")

			db, err := postgres.NewDB(&cfg.DB)
			if err != nil {
				log.Fatal().Err(err).Msg("DB connection failed")
			}

			ms := migrations.NewMigrationsService(db.GetConn(), cfg.Path)
			err = ms.ApplyPending()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to run DB migrations")
			}
		},
	}

	return migrationsCmd
}
