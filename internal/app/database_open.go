package app

import (
	"errors"

	"project/internal/config"

	neatdatabase "github.com/dracory/neat/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func databaseOpen(cfg config.ConfigInterface) (*neatdatabase.Database, error) {
	if cfg == nil {
		return nil, errors.New("databaseOpen: cfg is nil")
	}

	neatCfg := config.DatabaseNeatConfig(cfg)
	return neatdatabase.New(neatCfg)
}
