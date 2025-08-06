package repositories

import (
	"context"
	"days/internal/domain/entities"
	"days/internal/domain/valueobjects"
)

// UserRepository defines the contract for user persistence
type UserRepository interface {
	// Save persists a user entity
	Save(ctx context.Context, user *entities.User) error

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id *valueobjects.UserID) (*entities.User, error)

	// FindByEmail retrieves a user by email
	FindByEmail(ctx context.Context, email *valueobjects.Email) (*entities.User, error)

	// FindByUsername retrieves a user by username
	FindByUsername(ctx context.Context, username string) (*entities.User, error)

	// FindAll retrieves all users
	FindAll(ctx context.Context) ([]*entities.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error

	// Delete removes a user by ID
	Delete(ctx context.Context, id *valueobjects.UserID) error

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email *valueobjects.Email) (bool, error)

	// ExistsByUsername checks if a user with the given username exists
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
