package repositories

import (
	"context"
	"days/internal/domain/entities"
	"days/internal/domain/valueobjects"
)

// CalendarEntryRepository defines the contract for calendar entry persistence
type CalendarEntryRepository interface {
	// Save persists a calendar entry entity
	Save(ctx context.Context, entry *entities.CalendarEntry) error

	// FindByID retrieves a calendar entry by ID
	FindByID(ctx context.Context, id *valueobjects.UserID) (*entities.CalendarEntry, error)

	// FindByUserID retrieves all calendar entries for a user
	FindByUserID(ctx context.Context, userID *valueobjects.UserID) ([]*entities.CalendarEntry, error)

	// FindByUserIDAndDate retrieves a calendar entry for a specific user and date
	FindByUserIDAndDate(ctx context.Context, userID *valueobjects.UserID, date *valueobjects.CalendarDate) (*entities.CalendarEntry, error)

	// FindByUserIDAndDateRange retrieves calendar entries for a user within a date range
	FindByUserIDAndDateRange(ctx context.Context, userID *valueobjects.UserID, startDate, endDate *valueobjects.CalendarDate) ([]*entities.CalendarEntry, error)

	// FindByUserIDAndColorSetting retrieves calendar entries for a user with a specific color setting
	FindByUserIDAndColorSetting(ctx context.Context, userID *valueobjects.UserID, colorSettingID *valueobjects.UserID) ([]*entities.CalendarEntry, error)

	// Update updates an existing calendar entry
	Update(ctx context.Context, entry *entities.CalendarEntry) error

	// Delete removes a calendar entry by ID
	Delete(ctx context.Context, id *valueobjects.UserID) error

	// DeleteByUserIDAndDate removes a calendar entry by user ID and date
	DeleteByUserIDAndDate(ctx context.Context, userID *valueobjects.UserID, date *valueobjects.CalendarDate) error

	// ExistsByUserIDAndDate checks if a calendar entry exists for a user and date
	ExistsByUserIDAndDate(ctx context.Context, userID *valueobjects.UserID, date *valueobjects.CalendarDate) (bool, error)
}
