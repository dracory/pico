package seeders

import (
	"errors"

	"project/database/models"
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

	seedUsers := []*models.User{
		models.NewUser().
			SetName("Admin User").
			SetEmail("admin@example.com").
			SetPassword("$2a$10$N9qo8uLOickgx2ZMRZoMy.Mrq4JfZQg7q5Iq2v8qK5q5Iq2v8qK5q").
			SetRole(models.USER_ROLE_ADMINISTRATOR).
			SetStatus(models.USER_STATUS_ACTIVE),
		models.NewUser().
			SetName("Test User").
			SetEmail("user@example.com").
			SetPassword("$2a$10$N9qo8uLOickgx2ZMRZoMy.Mrq4JfZQg7q5Iq2v8qK5q5Iq2v8qK5q").
			SetRole(models.USER_ROLE_USER).
			SetStatus(models.USER_STATUS_ACTIVE),
	}

	for _, u := range seedUsers {
		err := query.Model(u).FirstOrCreate(u, "email = ?", u.Email)
		if err != nil {
			return err
		}
	}

	return nil
}
