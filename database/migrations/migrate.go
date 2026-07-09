package migrations

import (
	"context"
	"errors"
	"fmt"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

// MigrateAll runs all registered migrations using the neat migrator.
// Migrations are tracked in a single migration_tracker table.
func MigrateAll(a app.AppInterface) error {
	if a == nil {
		return errors.New("app is nil")
	}

	neatDB := a.GetNeatDatabase()
	if neatDB == nil {
		return errors.New("neat database is nil")
	}

	m := migrator.NewMigrator(neatDB)
	m.SetTransactionsEnabled(false)

	migrations := getMigrations()
	if err := m.AddMigrations(migrations); err != nil {
		return fmt.Errorf("failed to add migrations: %w", err)
	}

	return m.Up(context.Background())
}

// getMigrations returns all registered migrations.
func getMigrations() []migrator.MigrationInterface {
	return []migrator.MigrationInterface{
		&CreateUsersTable{},
	}
}
