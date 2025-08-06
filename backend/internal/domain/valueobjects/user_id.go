package valueobjects

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// UserID represents a unique identifier for a user
type UserID struct {
	value string
}

// NewUserID creates a new UserID from a string
func NewUserID(id string) (*UserID, error) {
	if id == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	return &UserID{value: id}, nil
}

// GenerateUserID creates a new random UserID
func GenerateUserID() *UserID {
	return &UserID{value: uuid.New().String()}
}

// Value returns the string value of the UserID
func (uid *UserID) Value() string {
	return uid.value
}

// String implements the Stringer interface
func (uid *UserID) String() string {
	return uid.value
}

// Equals checks if two UserIDs are equal
func (uid *UserID) Equals(other *UserID) bool {
	if other == nil {
		return false
	}
	return uid.value == other.value
}
