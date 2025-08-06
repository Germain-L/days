package services

import (
	"context"
	"days/internal/domain/entities"
	"days/internal/domain/repositories"
	"days/internal/domain/valueobjects"
	"errors"
)

// UserDomainService provides business logic operations for users
type UserDomainService struct {
	userRepo repositories.UserRepository
}

// NewUserDomainService creates a new UserDomainService
func NewUserDomainService(userRepo repositories.UserRepository) *UserDomainService {
	return &UserDomainService{
		userRepo: userRepo,
	}
}

// ValidateUniqueConstraints ensures email and username are unique
func (s *UserDomainService) ValidateUniqueConstraints(ctx context.Context, email *valueobjects.Email, username string, excludeUserID *valueobjects.UserID) error {
	// Check email uniqueness
	emailExists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if emailExists {
		// If we're updating a user, check if the email belongs to the same user
		if excludeUserID != nil {
			existingUser, err := s.userRepo.FindByEmail(ctx, email)
			if err != nil {
				return err
			}
			if existingUser != nil && !existingUser.ID().Equals(excludeUserID) {
				return errors.New("email already exists")
			}
		} else {
			return errors.New("email already exists")
		}
	}

	// Check username uniqueness
	usernameExists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return err
	}
	if usernameExists {
		// If we're updating a user, check if the username belongs to the same user
		if excludeUserID != nil {
			existingUser, err := s.userRepo.FindByUsername(ctx, username)
			if err != nil {
				return err
			}
			if existingUser != nil && !existingUser.ID().Equals(excludeUserID) {
				return errors.New("username already exists")
			}
		} else {
			return errors.New("username already exists")
		}
	}

	return nil
}

// ValidateUserExists checks if a user exists by ID
func (s *UserDomainService) ValidateUserExists(ctx context.Context, userID *valueobjects.UserID) (*entities.User, error) {
	if userID == nil {
		return nil, errors.New("user ID is required")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
