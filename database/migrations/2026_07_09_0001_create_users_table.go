package migrations

import (
	"errors"

	contractsSchema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*CreateUsersTable)(nil)

// CreateUsersTable creates the users table with id, name, email, password, and timestamps.
type CreateUsersTable struct {
	migrator.BaseMigration
}

func (m *CreateUsersTable) Signature() string {
	return "2026_07_09_0001_create_users_table"
}

func (m *CreateUsersTable) Description() string {
	return "Create users table with id, name, email, password, role, status, and timestamps"
}

func (m *CreateUsersTable) Up() error {
	schema := m.GetSchema()
	if schema == nil {
		return errors.New("schema is nil")
	}

	return schema.Create("users", func(table contractsSchema.Blueprint) {
		table.BigIncrements("id")
		table.String("name")
		table.String("email")
		table.Unique("email")
		table.String("password")
		table.String("role").Default("user")
		table.String("status").Default("active")
		table.Timestamps()
		table.SoftDeletes()
	})
}

func (m *CreateUsersTable) Down() error {
	schema := m.GetSchema()
	if schema == nil {
		return errors.New("schema is nil")
	}

	return schema.DropIfExists("users")
}
