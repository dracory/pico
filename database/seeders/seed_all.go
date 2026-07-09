package seeders

import (
	"errors"

	"project/internal/app"

	contractsSeeder "github.com/dracory/neat/contracts/database/seeder"
)

// SeedAll runs all registered seeders.
func SeedAll(a app.AppInterface) error {
	if a == nil {
		return errors.New("app is nil")
	}

	neatDB := a.GetNeatDatabase()
	if neatDB == nil {
		return errors.New("neat database is nil")
	}

	seeders := []contractsSeeder.Seeder{
		NewUserSeeder(a),
	}

	return neatDB.Seed(seeders)
}
