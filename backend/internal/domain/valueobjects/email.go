package valueobjects

import (
	"errors"
	"regexp"
)

// Email represents a valid email address
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email value object
func NewEmail(email string) (*Email, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	return &Email{value: email}, nil
}

// Value returns the string value of the email
func (e *Email) Value() string {
	return e.value
}

// String implements the Stringer interface
func (e *Email) String() string {
	return e.value
}

// Equals checks if two emails are equal
func (e *Email) Equals(other *Email) bool {
	if other == nil {
		return false
	}
	return e.value == other.value
}
