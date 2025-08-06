package usecases

import (
	"context"
	"days/internal/application/dto"
	"days/internal/domain/entities"
	"days/internal/domain/repositories"
	"days/internal/domain/services"
	"days/internal/domain/valueobjects"
	"errors"
)

// CalendarEntryUseCase handles calendar entry-related business operations
type CalendarEntryUseCase struct {
	calendarRepo      repositories.CalendarEntryRepository
	colorSettingRepo  repositories.ColorSettingRepository
	calendarDomainSvc *services.CalendarDomainService
}

// NewCalendarEntryUseCase creates a new CalendarEntryUseCase
func NewCalendarEntryUseCase(
	calendarRepo repositories.CalendarEntryRepository,
	colorSettingRepo repositories.ColorSettingRepository,
	calendarDomainSvc *services.CalendarDomainService,
) *CalendarEntryUseCase {
	return &CalendarEntryUseCase{
		calendarRepo:      calendarRepo,
		colorSettingRepo:  colorSettingRepo,
		calendarDomainSvc: calendarDomainSvc,
	}
}

// CreateCalendarEntry creates a new calendar entry
func (uc *CalendarEntryUseCase) CreateCalendarEntry(ctx context.Context, userID string, req *dto.CreateCalendarEntryRequest) (*dto.CalendarEntryResponse, error) {
	// Create value objects
	uid, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	date, err := valueobjects.NewCalendarDateFromString(req.Date)
	if err != nil {
		return nil, err
	}

	colorSettingID, err := valueobjects.NewUserID(req.ColorSettingID)
	if err != nil {
		return nil, err
	}

	// Validate business constraints
	err = uc.calendarDomainSvc.ValidateCalendarEntryConstraints(ctx, uid, date, colorSettingID, nil)
	if err != nil {
		return nil, err
	}

	// Create calendar entry entity
	entry, err := entities.NewCalendarEntry(uid, date, colorSettingID)
	if err != nil {
		return nil, err
	}

	// Set optional fields
	if req.Notes != nil {
		entry.UpdateNotes(req.Notes)
	}

	// Save calendar entry
	err = uc.calendarRepo.Save(ctx, entry)
	if err != nil {
		return nil, err
	}

	return uc.calendarEntryToResponse(ctx, entry)
}

// GetCalendarEntriesByUserID retrieves all calendar entries for a user
func (uc *CalendarEntryUseCase) GetCalendarEntriesByUserID(ctx context.Context, userID string, filters *dto.CalendarEntryFilters) (*dto.CalendarEntriesResponse, error) {
	uid, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	var entries []*entities.CalendarEntry

	// Apply filters
	if filters != nil && filters.StartDate != nil && filters.EndDate != nil {
		startDate, err := valueobjects.NewCalendarDateFromString(*filters.StartDate)
		if err != nil {
			return nil, err
		}
		endDate, err := valueobjects.NewCalendarDateFromString(*filters.EndDate)
		if err != nil {
			return nil, err
		}
		entries, err = uc.calendarRepo.FindByUserIDAndDateRange(ctx, uid, startDate, endDate)
		if err != nil {
			return nil, err
		}
	} else if filters != nil && filters.ColorSettingID != nil {
		colorSettingID, err := valueobjects.NewUserID(*filters.ColorSettingID)
		if err != nil {
			return nil, err
		}
		entries, err = uc.calendarRepo.FindByUserIDAndColorSetting(ctx, uid, colorSettingID)
		if err != nil {
			return nil, err
		}
	} else {
		entries, err = uc.calendarRepo.FindByUserID(ctx, uid)
		if err != nil {
			return nil, err
		}
	}

	responses := make([]dto.CalendarEntryResponse, len(entries))
	for i, entry := range entries {
		response, err := uc.calendarEntryToResponse(ctx, entry)
		if err != nil {
			return nil, err
		}
		responses[i] = *response
	}

	result := &dto.CalendarEntriesResponse{
		Data:  responses,
		Total: len(responses),
	}

	if filters != nil {
		result.Filters = *filters
	}

	return result, nil
}

// GetCalendarEntryByDate retrieves a calendar entry for a specific date
func (uc *CalendarEntryUseCase) GetCalendarEntryByDate(ctx context.Context, userID, dateStr string) (*dto.CalendarEntryResponse, error) {
	uid, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	date, err := valueobjects.NewCalendarDateFromString(dateStr)
	if err != nil {
		return nil, err
	}

	entry, err := uc.calendarRepo.FindByUserIDAndDate(ctx, uid, date)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.New("calendar entry not found for this date")
	}

	return uc.calendarEntryToResponse(ctx, entry)
}

// UpdateCalendarEntry updates an existing calendar entry
func (uc *CalendarEntryUseCase) UpdateCalendarEntry(ctx context.Context, userID, dateStr string, req *dto.UpdateCalendarEntryRequest) (*dto.CalendarEntryResponse, error) {
	uid, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	date, err := valueobjects.NewCalendarDateFromString(dateStr)
	if err != nil {
		return nil, err
	}

	// Get existing entry
	entry, err := uc.calendarRepo.FindByUserIDAndDate(ctx, uid, date)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.New("calendar entry not found for this date")
	}

	// Update color setting if provided
	if req.ColorSettingID != nil {
		colorSettingID, err := valueobjects.NewUserID(*req.ColorSettingID)
		if err != nil {
			return nil, err
		}

		// Validate business constraints
		err = uc.calendarDomainSvc.ValidateCalendarEntryConstraints(ctx, uid, date, colorSettingID, entry.ID())
		if err != nil {
			return nil, err
		}

		err = entry.UpdateColorSetting(colorSettingID)
		if err != nil {
			return nil, err
		}
	}

	// Update notes if provided
	if req.Notes != nil {
		entry.UpdateNotes(req.Notes)
	}

	// Save updated entry
	err = uc.calendarRepo.Update(ctx, entry)
	if err != nil {
		return nil, err
	}

	return uc.calendarEntryToResponse(ctx, entry)
}

// DeleteCalendarEntry deletes a calendar entry by date
func (uc *CalendarEntryUseCase) DeleteCalendarEntry(ctx context.Context, userID, dateStr string) error {
	uid, err := valueobjects.NewUserID(userID)
	if err != nil {
		return err
	}

	date, err := valueobjects.NewCalendarDateFromString(dateStr)
	if err != nil {
		return err
	}

	// Check if entry exists
	entry, err := uc.calendarRepo.FindByUserIDAndDate(ctx, uid, date)
	if err != nil {
		return err
	}
	if entry == nil {
		return errors.New("calendar entry not found for this date")
	}

	return uc.calendarRepo.DeleteByUserIDAndDate(ctx, uid, date)
}

// calendarEntryToResponse converts a calendar entry entity to a response DTO
func (uc *CalendarEntryUseCase) calendarEntryToResponse(ctx context.Context, entry *entities.CalendarEntry) (*dto.CalendarEntryResponse, error) {
	response := &dto.CalendarEntryResponse{
		ID:        entry.ID().Value(),
		UserID:    entry.UserID().Value(),
		Date:      entry.Date().String(),
		Notes:     entry.Notes(),
		CreatedAt: entry.CreatedAt(),
		UpdatedAt: entry.UpdatedAt(),
	}

	// Load color setting
	colorSetting, err := uc.colorSettingRepo.FindByID(ctx, entry.ColorSettingID())
	if err != nil {
		return nil, err
	}
	if colorSetting != nil {
		response.ColorSetting = &dto.ColorSettingResponse{
			ID:          colorSetting.ID().Value(),
			UserID:      colorSetting.UserID().Value(),
			Name:        colorSetting.Name(),
			HexColor:    colorSetting.HexColor().Value(),
			Description: colorSetting.Description(),
			IsDefault:   colorSetting.IsDefault(),
			SortOrder:   colorSetting.SortOrder(),
			CreatedAt:   colorSetting.CreatedAt(),
			UpdatedAt:   colorSetting.UpdatedAt(),
		}
	}

	return response, nil
}
