package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"days/internal/db"

	"github.com/google/uuid"
)

var (
	ErrDayEntryNotFound     = errors.New("day entry not found")
	ErrInvalidDate          = errors.New("invalid date format")
	ErrDayEntryExists       = errors.New("day entry already exists for this date")
	ErrUnauthorizedDayEntry = errors.New("not authorized to access this day entry")
)

type DayEntryService struct {
	queries             *db.Queries
	calendarService     *CalendarService
	colorMeaningService *ColorMeaningService
}

type CreateDayEntryRequest struct {
	Date           string    `json:"date"` // YYYY-MM-DD format
	ColorMeaningID uuid.UUID `json:"color_meaning_id"`
	Notes          *string   `json:"notes,omitempty"`
}

type UpdateDayEntryRequest struct {
	ColorMeaningID uuid.UUID `json:"color_meaning_id"`
	Notes          *string   `json:"notes,omitempty"`
}

type DayEntryResponse struct {
	ID             uuid.UUID `json:"id"`
	CalendarID     uuid.UUID `json:"calendar_id"`
	Date           string    `json:"date"`
	ColorMeaningID uuid.UUID `json:"color_meaning_id"`
	ColorHex       string    `json:"color_hex"`
	Meaning        string    `json:"meaning"`
	Notes          *string   `json:"notes,omitempty"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

type DateRangeRequest struct {
	StartDate string `json:"start_date"` // YYYY-MM-DD format
	EndDate   string `json:"end_date"`   // YYYY-MM-DD format
}

func NewDayEntryService(queries *db.Queries, calendarService *CalendarService, colorMeaningService *ColorMeaningService) *DayEntryService {
	return &DayEntryService{
		queries:             queries,
		calendarService:     calendarService,
		colorMeaningService: colorMeaningService,
	}
}

// CreateDayEntry creates a new day entry for a calendar
func (s *DayEntryService) CreateDayEntry(ctx context.Context, userID, calendarID uuid.UUID, req CreateDayEntryRequest) (*DayEntryResponse, error) {
	// Validate date
	date, err := s.parseDate(req.Date)
	if err != nil {
		return nil, err
	}

	// Check user owns the calendar
	_, err = s.calendarService.GetCalendarByID(ctx, userID, calendarID)
	if err != nil {
		return nil, err
	}

	// Check color meaning exists and belongs to this calendar
	colorMeaning, err := s.colorMeaningService.GetColorMeaningByID(ctx, userID, req.ColorMeaningID)
	if err != nil {
		return nil, fmt.Errorf("invalid color meaning: %w", err)
	}
	if colorMeaning.CalendarID != calendarID {
		return nil, errors.New("color meaning does not belong to this calendar")
	}

	// Check if day entry already exists for this date
	_, err = s.queries.GetDayEntryByCalendarAndDate(ctx, db.GetDayEntryByCalendarAndDateParams{
		CalendarID: calendarID,
		Date:       date,
	})
	if err == nil {
		return nil, ErrDayEntryExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing day entry: %w", err)
	}

	// Prepare notes
	var notes sql.NullString
	if req.Notes != nil {
		notes = sql.NullString{String: strings.TrimSpace(*req.Notes), Valid: true}
	}

	// Create day entry
	dayEntry, err := s.queries.CreateDayEntry(ctx, db.CreateDayEntryParams{
		CalendarID:     calendarID,
		Date:           date,
		ColorMeaningID: req.ColorMeaningID,
		Notes:          notes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create day entry: %w", err)
	}

	// Get the full day entry with color meaning details
	return s.getDayEntryWithColorMeaning(ctx, dayEntry.CalendarID, dayEntry.Date)
}

// GetDayEntriesByCalendarID retrieves all day entries for a calendar
func (s *DayEntryService) GetDayEntriesByCalendarID(ctx context.Context, userID, calendarID uuid.UUID) ([]*DayEntryResponse, error) {
	// Check user owns the calendar
	_, err := s.calendarService.GetCalendarByID(ctx, userID, calendarID)
	if err != nil {
		return nil, err
	}

	dayEntries, err := s.queries.GetDayEntriesByCalendarID(ctx, calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get day entries: %w", err)
	}

	var responses []*DayEntryResponse
	for _, de := range dayEntries {
		responses = append(responses, s.toDayEntryResponse(de))
	}

	return responses, nil
}

// GetDayEntriesByDateRange retrieves all day entries for a user within a date range
func (s *DayEntryService) GetDayEntriesByDateRange(ctx context.Context, userID uuid.UUID, req DateRangeRequest) ([]*DayEntryResponse, error) {
	startDate, err := s.parseDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := s.parseDate(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	if endDate.Before(startDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	dayEntries, err := s.queries.GetDayEntriesByDateRange(ctx, db.GetDayEntriesByDateRangeParams{
		UserID: userID,
		Date:   startDate,
		Date_2: endDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get day entries by date range: %w", err)
	}

	var responses []*DayEntryResponse
	for _, de := range dayEntries {
		response := &DayEntryResponse{
			ID:             de.ID,
			CalendarID:     de.CalendarID,
			Date:           de.Date.Format("2006-01-02"),
			ColorMeaningID: de.ColorMeaningID,
			ColorHex:       de.ColorHex,
			Meaning:        de.Meaning,
		}

		if de.Notes.Valid {
			response.Notes = &de.Notes.String
		}

		if de.CreatedAt.Valid {
			response.CreatedAt = de.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
		}
		if de.UpdatedAt.Valid {
			response.UpdatedAt = de.UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// GetDayEntryByCalendarAndDate retrieves a specific day entry
func (s *DayEntryService) GetDayEntryByCalendarAndDate(ctx context.Context, userID, calendarID uuid.UUID, dateStr string) (*DayEntryResponse, error) {
	// Validate date
	date, err := s.parseDate(dateStr)
	if err != nil {
		return nil, err
	}

	// Check user owns the calendar
	_, err = s.calendarService.GetCalendarByID(ctx, userID, calendarID)
	if err != nil {
		return nil, err
	}

	return s.getDayEntryWithColorMeaning(ctx, calendarID, date)
}

// UpdateDayEntry updates an existing day entry
func (s *DayEntryService) UpdateDayEntry(ctx context.Context, userID, calendarID uuid.UUID, dateStr string, req UpdateDayEntryRequest) (*DayEntryResponse, error) {
	// Validate date
	date, err := s.parseDate(dateStr)
	if err != nil {
		return nil, err
	}

	// Check user owns the calendar
	_, err = s.calendarService.GetCalendarByID(ctx, userID, calendarID)
	if err != nil {
		return nil, err
	}

	// Check day entry exists
	_, err = s.getDayEntryWithColorMeaning(ctx, calendarID, date)
	if err != nil {
		return nil, err
	}

	// Check color meaning exists and belongs to this calendar
	colorMeaning, err := s.colorMeaningService.GetColorMeaningByID(ctx, userID, req.ColorMeaningID)
	if err != nil {
		return nil, fmt.Errorf("invalid color meaning: %w", err)
	}
	if colorMeaning.CalendarID != calendarID {
		return nil, errors.New("color meaning does not belong to this calendar")
	}

	// Prepare notes
	var notes sql.NullString
	if req.Notes != nil {
		notes = sql.NullString{String: strings.TrimSpace(*req.Notes), Valid: true}
	}

	// Update day entry
	_, err = s.queries.UpdateDayEntry(ctx, db.UpdateDayEntryParams{
		CalendarID:     calendarID,
		ColorMeaningID: req.ColorMeaningID,
		Notes:          notes,
		Date:           date,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update day entry: %w", err)
	}

	// Return updated day entry
	return s.getDayEntryWithColorMeaning(ctx, calendarID, date)
}

// DeleteDayEntry deletes a day entry
func (s *DayEntryService) DeleteDayEntry(ctx context.Context, userID, calendarID uuid.UUID, dateStr string) error {
	// Validate date
	date, err := s.parseDate(dateStr)
	if err != nil {
		return err
	}

	// Check user owns the calendar
	_, err = s.calendarService.GetCalendarByID(ctx, userID, calendarID)
	if err != nil {
		return err
	}

	// Check day entry exists
	_, err = s.getDayEntryWithColorMeaning(ctx, calendarID, date)
	if err != nil {
		return err
	}

	// Delete day entry
	err = s.queries.DeleteDayEntry(ctx, db.DeleteDayEntryParams{
		CalendarID: calendarID,
		Date:       date,
	})
	if err != nil {
		return fmt.Errorf("failed to delete day entry: %w", err)
	}

	return nil
}

// Helper methods

func (s *DayEntryService) parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, ErrInvalidDate
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, ErrInvalidDate
	}

	return date, nil
}

func (s *DayEntryService) getDayEntryWithColorMeaning(ctx context.Context, calendarID uuid.UUID, date time.Time) (*DayEntryResponse, error) {
	dayEntry, err := s.queries.GetDayEntryByCalendarAndDate(ctx, db.GetDayEntryByCalendarAndDateParams{
		CalendarID: calendarID,
		Date:       date,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDayEntryNotFound
		}
		return nil, fmt.Errorf("failed to get day entry: %w", err)
	}

	return s.toDayEntryResponse(dayEntry), nil
}

func (s *DayEntryService) toDayEntryResponse(de interface{}) *DayEntryResponse {
	switch entry := de.(type) {
	case db.GetDayEntriesByCalendarIDRow:
		response := &DayEntryResponse{
			ID:             entry.ID,
			CalendarID:     entry.CalendarID,
			Date:           entry.Date.Format("2006-01-02"),
			ColorMeaningID: entry.ColorMeaningID,
			ColorHex:       entry.ColorHex,
			Meaning:        entry.Meaning,
		}

		if entry.Notes.Valid {
			response.Notes = &entry.Notes.String
		}

		if entry.CreatedAt.Valid {
			response.CreatedAt = entry.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
		}
		if entry.UpdatedAt.Valid {
			response.UpdatedAt = entry.UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
		}

		return response

	case db.GetDayEntryByCalendarAndDateRow:
		response := &DayEntryResponse{
			ID:             entry.ID,
			CalendarID:     entry.CalendarID,
			Date:           entry.Date.Format("2006-01-02"),
			ColorMeaningID: entry.ColorMeaningID,
			ColorHex:       entry.ColorHex,
			Meaning:        entry.Meaning,
		}

		if entry.Notes.Valid {
			response.Notes = &entry.Notes.String
		}

		if entry.CreatedAt.Valid {
			response.CreatedAt = entry.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
		}
		if entry.UpdatedAt.Valid {
			response.UpdatedAt = entry.UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
		}

		return response

	default:
		// This should never happen, but return a safe response
		return &DayEntryResponse{}
	}
}
