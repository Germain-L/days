package repositories

import (
	"context"
	"days/internal/domain/entities"
	"days/internal/domain/valueobjects"
)

// ColorSettingRepository defines the contract for color setting persistence
type ColorSettingRepository interface {
	// Save persists a color setting entity
	Save(ctx context.Context, colorSetting *entities.ColorSetting) error

	// FindByID retrieves a color setting by ID
	FindByID(ctx context.Context, id *valueobjects.UserID) (*entities.ColorSetting, error)

	// FindByUserID retrieves all color settings for a user
	FindByUserID(ctx context.Context, userID *valueobjects.UserID) ([]*entities.ColorSetting, error)

	// Update updates an existing color setting
	Update(ctx context.Context, colorSetting *entities.ColorSetting) error

	// Delete removes a color setting by ID
	Delete(ctx context.Context, id *valueobjects.UserID) error

	// FindByUserIDAndName finds a color setting by user ID and name
	FindByUserIDAndName(ctx context.Context, userID *valueobjects.UserID, name string) (*entities.ColorSetting, error)
}
