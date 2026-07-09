package app

import (
	"database/sql"
	"log/slog"

	"project/internal/config"

	neatdatabase "github.com/dracory/neat/database"
)

type AppInterface interface {
	Close() error

	GetLogger() *slog.Logger
	SetLogger(l *slog.Logger)

	GetConfig() config.ConfigInterface
	SetConfig(c config.ConfigInterface)

	GetDatabase() *sql.DB
	SetDatabase(db *sql.DB)

	GetNeatDatabase() *neatdatabase.Database
	SetNeatDatabase(db *neatdatabase.Database)
}
