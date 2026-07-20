package models

import (
	"time"

	"github.com/dracory/neat"
	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/support/uid"
)

const (
	UserTableName           = "user"
	USER_ROLE_ADMINISTRATOR = "administrator"
	USER_ROLE_USER          = "user"
	USER_STATUS_ACTIVE      = "active"
	USER_STATUS_INACTIVE    = "inactive"
)

// User represents a platform user. Embeds only orm.ShortID (string ID);
// datetime fields are stored as ISO 8601 strings per docs/domain/model.md §0.
// The original orm.Model / orm.SoftDeletes / orm.Timestamps embeddings are
// dropped because they use ID uint, sql.NullTime, and time.Time respectively,
// which conflict with the string-ID / no-NULL / datetime-as-string rules.
type User struct {
	orm.ShortID
	Status         string `json:"status" db:"status"`
	FirstName      string `json:"first_name" db:"first_name"`
	MiddleNames    string `json:"middle_names" db:"middle_names"`
	LastName       string `json:"last_name" db:"last_name"`
	Email          string `json:"email" db:"email"`
	Password       string `json:"-" db:"password"`
	Role           string `json:"role" db:"role"`
	CreatedAt      string `json:"created_at" db:"created_at"`
	UpdatedAt      string `json:"updated_at" db:"updated_at"`
	SoftDeletedAt  string `json:"soft_deleted_at" db:"soft_deleted_at"`
}

// TableName returns the database table name for User.
func (User) TableName() string {
	return UserTableName
}

// NewUser creates a new User instance with a generated short ID, active
// status, user role, and datetime fields initialized to current UTC
// (CreatedAt/UpdatedAt) or the neat.NullDateTime sentinel (SoftDeletedAt).
func NewUser() *User {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	return &User{
		ShortID:       orm.ShortID{ID: uid.GenerateShortID()},
		Status:        USER_STATUS_ACTIVE,
		Role:          USER_ROLE_USER,
		CreatedAt:     now,
		UpdatedAt:     now,
		SoftDeletedAt: neat.NullDateTime,
	}
}

// SetID sets the ID and returns the pointer for chaining.
func (u *User) SetID(id string) *User {
	u.ID = id
	return u
}

// SetStatus sets the status and returns the pointer for chaining.
func (u *User) SetStatus(status string) *User {
	u.Status = status
	return u
}

// SetFirstName sets the first name and returns the pointer for chaining.
func (u *User) SetFirstName(firstName string) *User {
	u.FirstName = firstName
	return u
}

// SetMiddleNames sets the middle names and returns the pointer for chaining.
func (u *User) SetMiddleNames(middleNames string) *User {
	u.MiddleNames = middleNames
	return u
}

// SetLastName sets the last name and returns the pointer for chaining.
func (u *User) SetLastName(lastName string) *User {
	u.LastName = lastName
	return u
}

// SetEmail sets the email and returns the pointer for chaining.
func (u *User) SetEmail(email string) *User {
	u.Email = email
	return u
}

// SetPassword sets the password and returns the pointer for chaining.
func (u *User) SetPassword(password string) *User {
	u.Password = password
	return u
}

// SetRole sets the role and returns the pointer for chaining.
func (u *User) SetRole(role string) *User {
	u.Role = role
	return u
}

// SetCreatedAt sets the created_at datetime string and returns the pointer for chaining.
func (u *User) SetCreatedAt(s string) *User {
	u.CreatedAt = s
	return u
}

// SetUpdatedAt sets the updated_at datetime string and returns the pointer for chaining.
func (u *User) SetUpdatedAt(s string) *User {
	u.UpdatedAt = s
	return u
}
