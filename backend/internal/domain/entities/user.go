package entities

import (
	"days/internal/domain/valueobjects"
	"errors"
	"time"
)

// User represents a user entity in the domain
type User struct {
	id        *valueobjects.UserID
	email     *valueobjects.Email
	username  string
	firstName *string
	lastName  *string
	timezone  string
	createdAt time.Time
	updatedAt time.Time
}

// NewUser creates a new User entity
func NewUser(email *valueobjects.Email, username string) (*User, error) {
	if email == nil {
		return nil, errors.New("email is required")
	}
	if username == "" {
		return nil, errors.New("username is required")
	}
	if len(username) < 3 || len(username) > 50 {
		return nil, errors.New("username must be between 3 and 50 characters")
	}

	now := time.Now()
	return &User{
		id:        valueobjects.GenerateUserID(),
		email:     email,
		username:  username,
		timezone:  "UTC",
		createdAt: now,
		updatedAt: now,
	}, nil
}

// LoadUser loads an existing user with all data
func LoadUser(
	id *valueobjects.UserID,
	email *valueobjects.Email,
	username string,
	firstName *string,
	lastName *string,
	timezone string,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:        id,
		email:     email,
		username:  username,
		firstName: firstName,
		lastName:  lastName,
		timezone:  timezone,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID returns the user's ID
func (u *User) ID() *valueobjects.UserID {
	return u.id
}

// Email returns the user's email
func (u *User) Email() *valueobjects.Email {
	return u.email
}

// Username returns the user's username
func (u *User) Username() string {
	return u.username
}

// FirstName returns the user's first name
func (u *User) FirstName() *string {
	return u.firstName
}

// LastName returns the user's last name
func (u *User) LastName() *string {
	return u.lastName
}

// Timezone returns the user's timezone
func (u *User) Timezone() string {
	return u.timezone
}

// CreatedAt returns when the user was created
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns when the user was last updated
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(email *valueobjects.Email) error {
	if email == nil {
		return errors.New("email is required")
	}
	u.email = email
	u.updatedAt = time.Now()
	return nil
}

// UpdateUsername updates the user's username
func (u *User) UpdateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if len(username) < 3 || len(username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	u.username = username
	u.updatedAt = time.Now()
	return nil
}

// UpdateFirstName updates the user's first name
func (u *User) UpdateFirstName(firstName *string) {
	u.firstName = firstName
	u.updatedAt = time.Now()
}

// UpdateLastName updates the user's last name
func (u *User) UpdateLastName(lastName *string) {
	u.lastName = lastName
	u.updatedAt = time.Now()
}

// UpdateTimezone updates the user's timezone
func (u *User) UpdateTimezone(timezone string) {
	if timezone == "" {
		timezone = "UTC"
	}
	u.timezone = timezone
	u.updatedAt = time.Now()
}
