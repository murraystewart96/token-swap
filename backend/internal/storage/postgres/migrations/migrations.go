package migrations

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

type Migrations struct {
	conn           *pgxpool.Pool
	migrationsPath string
}

func NewMigrationsService(conn *pgxpool.Pool, migrationsPath string) *Migrations {
	return &Migrations{conn: conn, migrationsPath: migrationsPath}
}

func (m *Migrations) ApplyPending() error {
	currentVersion, err := m.getCurrentVersion()
	if err != nil {
		return err
	}

	if err := goose.Up(stdlib.OpenDBFromPool(m.conn), m.migrationsPath); err != nil {
		return fmt.Errorf(
			"goose.Up : version: [%v] | [%w]",
			currentVersion, err,
		)
	}

	return nil
}

func (m *Migrations) CheckForPendingMigrations() error {
	currentVersion, err := m.getCurrentVersion()
	if err != nil {
		return err
	}

	latestPathVersion, err := m.getLatestMigrationFromPath()
	if err != nil {
		return err
	}

	if currentVersion < latestPathVersion {
		msg := "has unproccessed migrations in path"

		return fmt.Errorf("migrations: error: %s current: %d - latest: %d", msg, currentVersion, latestPathVersion) //nolint:err113
	}

	return nil
}

func (m *Migrations) getCurrentVersion() (int64, error) {
	current, err := goose.EnsureDBVersion(stdlib.OpenDBFromPool(m.conn))
	if err != nil {
		comment := "get current goose_db_version"

		return 0, fmt.Errorf("goose.EnsureDBVersion : [%s] | [%w]", comment, err)
	}

	return current, nil
}

func (m *Migrations) getLatestMigrationFromPath() (int64, error) {
	migrations, err := goose.CollectMigrations(m.migrationsPath, 0, goose.MaxVersion)
	if err != nil {
		return 0, fmt.Errorf("failed to get list of migrations: %w", err)
	}

	// The latest migration version should be the highest version
	return migrations[len(migrations)-1].Version, nil
}
