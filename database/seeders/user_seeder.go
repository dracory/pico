package seeders

import (
	"errors"

	"project/internal/app"

	contractsSeeder "github.com/dracory/neat/contracts/database/seeder"
)

var _ contractsSeeder.Seeder = (*UserSeeder)(nil)

// UserSeeder inserts default user records into the database.
type UserSeeder struct {
	app app.AppInterface
}

// NewUserSeeder creates a new UserSeeder instance.
func NewUserSeeder(a app.AppInterface) *UserSeeder {
	return &UserSeeder{app: a}
}

func (s *UserSeeder) Signature() string {
	return "2026_07_09_0001_user_seeder"
}

func (s *UserSeeder) Run() error {
	if s.app == nil {
		return errors.New("app is nil")
	}

	neatDB := s.app.GetNeatDatabase()
	if neatDB == nil {
		return errors.New("neat database is nil")
	}

	query := neatDB.Query()

	users := []map[string]any{
		{
			"name":     "Admin User",
			"email":    "admin@example.com",
			"password": "$2a$10$N9qo8uLOickgx2ZMRZoMy.Mrq4JfZQg7q5Iq2v8qK5q5Iq2v8qK5q",
			"role":     "administrator",
			"status":   "active",
		},
		{
			"name":     "Test User",
			"email":    "user@example.com",
			"password": "$2a$10$N9qo8uLOickgx2ZMRZoMy.Mrq4JfZQg7q5Iq2v8qK5q5Iq2v8qK5q",
			"role":     "user",
			"status":   "active",
		},
	}

	for _, user := range users {
		err := query.Table("users").UpdateOrInsert(
			map[string]any{"email": user["email"]},
			user,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
