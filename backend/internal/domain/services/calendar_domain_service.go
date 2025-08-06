package services

import (
	"context"
	"days/internal/domain/entities"
	"days/internal/domain/repositories"
	"days/internal/domain/valueobjects"
	"errors"
)

// CalendarDomainService provides business logic operations for calendar entries
type CalendarDomainService struct {
	calendarRepo     repositories.CalendarEntryRepository
	colorSettingRepo repositories.ColorSettingRepository
	userDomainSvc    *UserDomainService
}

// NewCalendarDomainService creates a new CalendarDomainService
func NewCalendarDomainService(
	calendarRepo repositories.CalendarEntryRepository,
	colorSettingRepo repositories.ColorSettingRepository,
	userDomainSvc *UserDomainService,
) *CalendarDomainService {
	return &CalendarDomainService{
		calendarRepo:     calendarRepo,
		colorSettingRepo: colorSettingRepo,
		userDomainSvc:    userDomainSvc,
	}
}

// ValidateCalendarEntryConstraints validates business rules for calendar entries
func (s *CalendarDomainService) ValidateCalendarEntryConstraints(
	ctx context.Context,
	userID *valueobjects.UserID,
	date *valueobjects.CalendarDate,
	colorSettingID *valueobjects.UserID,
	excludeEntryID *valueobjects.UserID,
) error {
	// Validate user exists
	_, err := s.userDomainSvc.ValidateUserExists(ctx, userID)
	if err != nil {
		return err
	}

	// Validate color setting exists and belongs to user
	colorSetting, err := s.colorSettingRepo.FindByID(ctx, colorSettingID)
	if err != nil {
		return err
	}
	if colorSetting == nil {
		return errors.New("color setting not found")
	}
	if !colorSetting.BelongsToUser(userID) {
		return errors.New("color setting does not belong to user")
	}

	// Check for duplicate entry on the same date (one entry per user per date)
	existingEntry, err := s.calendarRepo.FindByUserIDAndDate(ctx, userID, date)
	if err != nil {
		return err
	}
	if existingEntry != nil {
		// If we're updating an entry, make sure it's not a different entry for the same date
		if excludeEntryID == nil || !existingEntry.ID().Equals(excludeEntryID) {
			return errors.New("calendar entry already exists for this date")
		}
	}

	return nil
}

// ValidateColorSettingOwnership ensures a color setting belongs to the specified user
func (s *CalendarDomainService) ValidateColorSettingOwnership(
	ctx context.Context,
	userID *valueobjects.UserID,
	colorSettingID *valueobjects.UserID,
) (*entities.ColorSetting, error) {
	colorSetting, err := s.colorSettingRepo.FindByID(ctx, colorSettingID)
	if err != nil {
		return nil, err
	}
	if colorSetting == nil {
		return nil, errors.New("color setting not found")
	}
	if !colorSetting.BelongsToUser(userID) {
		return nil, errors.New("color setting does not belong to user")
	}
	return colorSetting, nil
}

// ValidateCalendarEntryOwnership ensures a calendar entry belongs to the specified user
func (s *CalendarDomainService) ValidateCalendarEntryOwnership(
	ctx context.Context,
	userID *valueobjects.UserID,
	entryID *valueobjects.UserID,
) (*entities.CalendarEntry, error) {
	entry, err := s.calendarRepo.FindByID(ctx, entryID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.New("calendar entry not found")
	}
	if !entry.BelongsToUser(userID) {
		return nil, errors.New("calendar entry does not belong to user")
	}
	return entry, nil
}
