package usecases

import (
	"context"
	"days/internal/application/dto"
	"days/internal/domain/entities"
	"days/internal/domain/repositories"
	"days/internal/domain/services"
	"days/internal/domain/valueobjects"
)

// UserUseCase handles user-related business operations
type UserUseCase struct {
	userRepo      repositories.UserRepository
	userDomainSvc *services.UserDomainService
}

// NewUserUseCase creates a new UserUseCase
func NewUserUseCase(
	userRepo repositories.UserRepository,
	userDomainSvc *services.UserDomainService,
) *UserUseCase {
	return &UserUseCase{
		userRepo:      userRepo,
		userDomainSvc: userDomainSvc,
	}
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Create value objects
	email, err := valueobjects.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Validate unique constraints
	err = uc.userDomainSvc.ValidateUniqueConstraints(ctx, email, req.Username, nil)
	if err != nil {
		return nil, err
	}

	// Create user entity
	user, err := entities.NewUser(email, req.Username)
	if err != nil {
		return nil, err
	}

	// Set optional fields
	if req.FirstName != nil {
		user.UpdateFirstName(req.FirstName)
	}
	if req.LastName != nil {
		user.UpdateLastName(req.LastName)
	}
	if req.Timezone != nil {
		user.UpdateTimezone(*req.Timezone)
	}

	// Save user
	err = uc.userRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	return uc.userToResponse(user), nil
}

// GetUser retrieves a user by ID
func (uc *UserUseCase) GetUser(ctx context.Context, userID string) (*dto.UserResponse, error) {
	id, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userDomainSvc.ValidateUserExists(ctx, id)
	if err != nil {
		return nil, err
	}

	return uc.userToResponse(user), nil
}

// GetAllUsers retrieves all users
func (uc *UserUseCase) GetAllUsers(ctx context.Context) (*dto.UsersResponse, error) {
	users, err := uc.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *uc.userToResponse(user)
	}

	return &dto.UsersResponse{
		Data:  responses,
		Total: len(responses),
	}, nil
}

// UpdateUser updates an existing user
func (uc *UserUseCase) UpdateUser(ctx context.Context, userID string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	id, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userDomainSvc.ValidateUserExists(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate unique constraints if email or username is being updated
	if req.Email != nil || req.Username != nil {
		email := user.Email()
		username := user.Username()

		if req.Email != nil {
			email, err = valueobjects.NewEmail(*req.Email)
			if err != nil {
				return nil, err
			}
		}
		if req.Username != nil {
			username = *req.Username
		}

		err = uc.userDomainSvc.ValidateUniqueConstraints(ctx, email, username, id)
		if err != nil {
			return nil, err
		}
	}

	// Update fields
	if req.Email != nil {
		email, err := valueobjects.NewEmail(*req.Email)
		if err != nil {
			return nil, err
		}
		err = user.UpdateEmail(email)
		if err != nil {
			return nil, err
		}
	}
	if req.Username != nil {
		err = user.UpdateUsername(*req.Username)
		if err != nil {
			return nil, err
		}
	}
	if req.FirstName != nil {
		user.UpdateFirstName(req.FirstName)
	}
	if req.LastName != nil {
		user.UpdateLastName(req.LastName)
	}
	if req.Timezone != nil {
		user.UpdateTimezone(*req.Timezone)
	}

	// Save updated user
	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return uc.userToResponse(user), nil
}

// DeleteUser deletes a user by ID
func (uc *UserUseCase) DeleteUser(ctx context.Context, userID string) error {
	id, err := valueobjects.NewUserID(userID)
	if err != nil {
		return err
	}

	_, err = uc.userDomainSvc.ValidateUserExists(ctx, id)
	if err != nil {
		return err
	}

	return uc.userRepo.Delete(ctx, id)
}

// userToResponse converts a user entity to a response DTO
func (uc *UserUseCase) userToResponse(user *entities.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID().Value(),
		Email:     user.Email().Value(),
		Username:  user.Username(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Timezone:  user.Timezone(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}
