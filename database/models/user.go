package models

import (
	"time"

	"github.com/dracory/neat/database/orm"
)

const (
	USER_ROLE_ADMINISTRATOR = "administrator"
	USER_ROLE_USER          = "user"
	USER_STATUS_ACTIVE      = "active"
	USER_STATUS_INACTIVE    = "inactive"
)

// User represents a user record in the database.
// It embeds orm.Model (ID + Timestamps) and orm.DeletedAt (Laravel-compatible soft deletes).
type User struct {
	orm.Model
	orm.SoftDeletes
	Name     string
	Email    string
	Password string
	Role     string
	Status   string
}

// TableName returns the database table name for User.
func (User) TableName() string {
	return "users"
}

// NewUser creates a new User instance with sensible defaults.
func NewUser() *User {
	return &User{
		Role:   USER_ROLE_USER,
		Status: USER_STATUS_ACTIVE,
	}
}

// SetName sets the name and returns the pointer for chaining.
func (u *User) SetName(name string) *User {
	u.Name = name
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

// SetStatus sets the status and returns the pointer for chaining.
func (u *User) SetStatus(status string) *User {
	u.Status = status
	return u
}

// SetCreatedAt sets the created_at timestamp and returns the pointer for chaining.
func (u *User) SetCreatedAt(t time.Time) *User {
	u.CreatedAt = t
	return u
}

// SetUpdatedAt sets the updated_at timestamp and returns the pointer for chaining.
func (u *User) SetUpdatedAt(t time.Time) *User {
	u.UpdatedAt = t
	return u
}
