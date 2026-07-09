package app

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"project/internal/config"

	neatdatabase "github.com/dracory/neat/database"
	"github.com/lmittmann/tint"
)

type appImplementation struct {
	cfg config.ConfigInterface

	neatDB *neatdatabase.Database
	db     *sql.DB

	logger *slog.Logger
}

var _ AppInterface = (*appImplementation)(nil)

func New(cfg config.ConfigInterface) (AppInterface, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}

	logger := slog.New(tint.NewHandler(os.Stdout, nil))

	neatDB, err := databaseOpen(cfg)
	if err != nil {
		return nil, err
	}

	db, err := neatDB.DB()
	if err != nil {
		return nil, err
	}

	a := &appImplementation{cfg: cfg}
	a.SetLogger(logger)
	a.SetNeatDatabase(neatDB)
	a.SetDatabase(db)

	return a, nil
}

func (r *appImplementation) Close() error {
	if r == nil {
		return nil
	}
	if r.neatDB == nil {
		return nil
	}
	err := r.neatDB.Close()
	r.neatDB = nil
	r.db = nil
	return err
}

func (r *appImplementation) GetConfig() config.ConfigInterface {
	if r == nil {
		return nil
	}
	return r.cfg
}

func (r *appImplementation) SetConfig(cfg config.ConfigInterface) {
	r.cfg = cfg
}

func (r *appImplementation) GetDatabase() *sql.DB {
	return r.db
}

func (r *appImplementation) SetDatabase(db *sql.DB) {
	r.db = db
}

func (r *appImplementation) GetNeatDatabase() *neatdatabase.Database {
	return r.neatDB
}

func (r *appImplementation) SetNeatDatabase(db *neatdatabase.Database) {
	r.neatDB = db
}

func (r *appImplementation) GetLogger() *slog.Logger {
	return r.logger
}

func (r *appImplementation) SetLogger(l *slog.Logger) {
	r.logger = l
}
