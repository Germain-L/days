package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"days/internal/db"

	"github.com/google/uuid"
)

var (
	ErrCalendarNotFound     = errors.New("calendar not found")
	ErrCalendarNameEmpty    = errors.New("calendar name cannot be empty")
	ErrCalendarNameExists   = errors.New("calendar with this name already exists")
	ErrUnauthorizedCalendar = errors.New("not authorized to access this calendar")
)

type CalendarService struct {
	queries *db.Queries
}

type CreateCalendarRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type UpdateCalendarRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type CalendarResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func NewCalendarService(queries *db.Queries) *CalendarService {
	return &CalendarService{
		queries: queries,
	}
}

// CreateCalendar creates a new calendar for a user
func (s *CalendarService) CreateCalendar(ctx context.Context, userID uuid.UUID, req CreateCalendarRequest) (*CalendarResponse, error) {
	// Validate input
	if err := s.validateCalendarName(req.Name); err != nil {
		return nil, err
	}

	// Check if user already has a calendar with this name
	existingCalendars, err := s.queries.GetCalendarsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing calendars: %w", err)
	}

	for _, calendar := range existingCalendars {
		if strings.EqualFold(calendar.Name, req.Name) {
			return nil, ErrCalendarNameExists
		}
	}

	// Prepare description
	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}

	// Create calendar
	calendar, err := s.queries.CreateCalendar(ctx, db.CreateCalendarParams{
		UserID:      userID,
		Name:        strings.TrimSpace(req.Name),
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar: %w", err)
	}

	return s.toCalendarResponse(calendar), nil
}

// GetCalendarsByUserID retrieves all calendars for a user
func (s *CalendarService) GetCalendarsByUserID(ctx context.Context, userID uuid.UUID) ([]*CalendarResponse, error) {
	calendars, err := s.queries.GetCalendarsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user calendars: %w", err)
	}

	var responses []*CalendarResponse
	for _, calendar := range calendars {
		responses = append(responses, s.toCalendarResponse(calendar))
	}

	return responses, nil
}

// GetCalendarByID retrieves a calendar by ID and verifies user ownership
func (s *CalendarService) GetCalendarByID(ctx context.Context, userID, calendarID uuid.UUID) (*CalendarResponse, error) {
	calendar, err := s.queries.GetCalendarByID(ctx, calendarID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCalendarNotFound
		}
		return nil, fmt.Errorf("failed to get calendar: %w", err)
	}

	// Check ownership
	if calendar.UserID != userID {
		return nil, ErrUnauthorizedCalendar
	}

	return s.toCalendarResponse(calendar), nil
}

// UpdateCalendar updates a calendar's name and description
func (s *CalendarService) UpdateCalendar(ctx context.Context, userID, calendarID uuid.UUID, req UpdateCalendarRequest) (*CalendarResponse, error) {
	// Validate input
	if err := s.validateCalendarName(req.Name); err != nil {
		return nil, err
	}

	// Check calendar exists and user owns it
	_, err := s.GetCalendarByID(ctx, userID, calendarID)
	if err != nil {
		return nil, err
	}

	// Check if user already has another calendar with this name
	userCalendars, err := s.queries.GetCalendarsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing calendars: %w", err)
	}

	for _, calendar := range userCalendars {
		if calendar.ID != calendarID && strings.EqualFold(calendar.Name, req.Name) {
			return nil, ErrCalendarNameExists
		}
	}

	// Prepare description
	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}

	// Update calendar
	updatedCalendar, err := s.queries.UpdateCalendar(ctx, db.UpdateCalendarParams{
		ID:          calendarID,
		Name:        strings.TrimSpace(req.Name),
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update calendar: %w", err)
	}

	// Verify the update didn't change ownership (shouldn't happen, but safety check)
	if updatedCalendar.UserID != userID {
		return nil, ErrUnauthorizedCalendar
	}

	return s.toCalendarResponse(updatedCalendar), nil
}

// DeleteCalendar deletes a calendar and all associated data
func (s *CalendarService) DeleteCalendar(ctx context.Context, userID, calendarID uuid.UUID) error {
	// Check calendar exists and user owns it
	_, err := s.GetCalendarByID(ctx, userID, calendarID)
	if err != nil {
		return err
	}

	// Delete calendar (cascades to color_meanings and day_entries)
	err = s.queries.DeleteCalendar(ctx, calendarID)
	if err != nil {
		return fmt.Errorf("failed to delete calendar: %w", err)
	}

	return nil
}

// Helper methods

func (s *CalendarService) validateCalendarName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrCalendarNameEmpty
	}
	if len(name) > 100 {
		return errors.New("calendar name cannot exceed 100 characters")
	}
	return nil
}

func (s *CalendarService) toCalendarResponse(calendar db.Calendar) *CalendarResponse {
	var description *string
	if calendar.Description.Valid {
		description = &calendar.Description.String
	}

	var createdAt, updatedAt string
	if calendar.CreatedAt.Valid {
		createdAt = calendar.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
	}
	if calendar.UpdatedAt.Valid {
		updatedAt = calendar.UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
	}

	return &CalendarResponse{
		ID:          calendar.ID,
		UserID:      calendar.UserID,
		Name:        calendar.Name,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
