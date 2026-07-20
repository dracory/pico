package migrations

import (
	"errors"

	"project/database/models"

	contractsSchema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dracory/neat"
	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*CreateUsersTable)(nil)

// CreateUsersTable creates the user table (singular) with varchar(21) string IDs,
// datetime columns, NOT NULL constraints, and soft_deleted_at with the
// neat.NullDateTime sentinel. Rewritten in place from the original Pico starter
// `users` migration per docs/domain/model.md §4.
type CreateUsersTable struct {
	migrator.BaseMigration
}

func (m *CreateUsersTable) Signature() string {
	return "2026_07_09_0001_create_users_table"
}

func (m *CreateUsersTable) Description() string {
	return "Create user table (singular) with string ID, split name fields, datetime columns, and soft_deleted_at sentinel"
}

func (m *CreateUsersTable) Up() error {
	schema := m.GetSchema()
	if schema == nil {
		return errors.New("schema is nil")
	}

	return schema.Create(models.UserTableName, func(table contractsSchema.Blueprint) {
		table.String("id", 21)
		table.Primary("id")
		table.String("status", 20).Default("active")
		table.String("first_name", 50)
		table.String("middle_names", 50).Default("")
		table.String("last_name", 50)
		table.String("email", 100)
		table.Unique("email")
		table.String("password", 255)
		table.String("role", 20).Default("user")
		table.DateTime("created_at").Default(neat.NullDateTime)
		table.DateTime("updated_at").Default(neat.NullDateTime)
		table.DateTime("soft_deleted_at").Default(neat.NullDateTime)
		table.Index("status")
	})
}

func (m *CreateUsersTable) Down() error {
	schema := m.GetSchema()
	if schema == nil {
		return errors.New("schema is nil")
	}

	return schema.DropIfExists(models.UserTableName)
}
