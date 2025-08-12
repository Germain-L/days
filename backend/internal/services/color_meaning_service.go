package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"days/internal/db"

	"github.com/google/uuid"
)

var (
	ErrColorMeaningNotFound     = errors.New("color meaning not found")
	ErrInvalidColorHex          = errors.New("invalid color hex format")
	ErrMeaningEmpty             = errors.New("meaning cannot be empty")
	ErrColorMeaningExists       = errors.New("color or meaning already exists for this calendar")
	ErrUnauthorizedColorMeaning = errors.New("not authorized to access this color meaning")
)

type ColorMeaningService struct {
	queries         *db.Queries
	calendarService *CalendarService
}

type CreateColorMeaningRequest struct {
	ColorHex string `json:"color_hex"`
	Meaning  string `json:"meaning"`
}

type UpdateColorMeaningRequest struct {
	ColorHex string `json:"color_hex"`
	Meaning  string `json:"meaning"`
}

type ColorMeaningResponse struct {
	ID         uuid.UUID `json:"id"`
	CalendarID uuid.UUID `json:"calendar_id"`
	ColorHex   string    `json:"color_hex"`
	Meaning    string    `json:"meaning"`
	CreatedAt  string    `json:"created_at"`
}

func NewColorMeaningService(queries *db.Queries, calendarService *CalendarService) *ColorMeaningService {
	return &ColorMeaningService{
		queries:         queries,
		calendarService: calendarService,
	}
}

// CreateColorMeaning creates a new color meaning for a calendar
func (s *ColorMeaningService) CreateColorMeaning(ctx context.Context, userID, calendarID uuid.UUID, req CreateColorMeaningRequest) (*ColorMeaningResponse, error) {
	// Validate input
	if err := s.validateColorHex(req.ColorHex); err != nil {
		return nil, err
	}
	if err := s.validateMeaning(req.Meaning); err != nil {
		return nil, err
	}

	// Check user owns the calendar
	_, err := s.calendarService.GetCalendarByID(ctx, userID, calendarID)
	if err != nil {
		return nil, err
	}

	// Check if color or meaning already exists for this calendar
	existingColorMeanings, err := s.queries.GetColorMeaningsByCalendarID(ctx, calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing color meanings: %w", err)
	}

	normalizedColorHex := s.normalizeColorHex(req.ColorHex)
	normalizedMeaning := strings.TrimSpace(req.Meaning)

	for _, cm := range existingColorMeanings {
		if s.normalizeColorHex(cm.ColorHex) == normalizedColorHex {
			return nil, errors.New("color already exists for this calendar")
		}
		if strings.EqualFold(cm.Meaning, normalizedMeaning) {
			return nil, errors.New("meaning already exists for this calendar")
		}
	}

	// Create color meaning
	colorMeaning, err := s.queries.CreateColorMeaning(ctx, db.CreateColorMeaningParams{
		CalendarID: calendarID,
		ColorHex:   normalizedColorHex,
		Meaning:    normalizedMeaning,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create color meaning: %w", err)
	}

	return s.toColorMeaningResponse(colorMeaning), nil
}

// GetColorMeaningsByCalendarID retrieves all color meanings for a calendar
func (s *ColorMeaningService) GetColorMeaningsByCalendarID(ctx context.Context, userID, calendarID uuid.UUID) ([]*ColorMeaningResponse, error) {
	// Check user owns the calendar
	_, err := s.calendarService.GetCalendarByID(ctx, userID, calendarID)
	if err != nil {
		return nil, err
	}

	colorMeanings, err := s.queries.GetColorMeaningsByCalendarID(ctx, calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get color meanings: %w", err)
	}

	var responses []*ColorMeaningResponse
	for _, cm := range colorMeanings {
		responses = append(responses, s.toColorMeaningResponse(cm))
	}

	return responses, nil
}

// GetColorMeaningByID retrieves a color meaning by ID and verifies user access
func (s *ColorMeaningService) GetColorMeaningByID(ctx context.Context, userID, colorMeaningID uuid.UUID) (*ColorMeaningResponse, error) {
	colorMeaning, err := s.queries.GetColorMeaningByID(ctx, colorMeaningID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrColorMeaningNotFound
		}
		return nil, fmt.Errorf("failed to get color meaning: %w", err)
	}

	// Check user owns the calendar
	_, err = s.calendarService.GetCalendarByID(ctx, userID, colorMeaning.CalendarID)
	if err != nil {
		if errors.Is(err, ErrUnauthorizedCalendar) {
			return nil, ErrUnauthorizedColorMeaning
		}
		return nil, err
	}

	return s.toColorMeaningResponse(colorMeaning), nil
}

// UpdateColorMeaning updates a color meaning
func (s *ColorMeaningService) UpdateColorMeaning(ctx context.Context, userID, colorMeaningID uuid.UUID, req UpdateColorMeaningRequest) (*ColorMeaningResponse, error) {
	// Validate input
	if err := s.validateColorHex(req.ColorHex); err != nil {
		return nil, err
	}
	if err := s.validateMeaning(req.Meaning); err != nil {
		return nil, err
	}

	// Get existing color meaning and verify access
	existingColorMeaning, err := s.GetColorMeaningByID(ctx, userID, colorMeaningID)
	if err != nil {
		return nil, err
	}

	// Check if new color or meaning conflicts with other color meanings in the same calendar
	calendarColorMeanings, err := s.queries.GetColorMeaningsByCalendarID(ctx, existingColorMeaning.CalendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing color meanings: %w", err)
	}

	normalizedColorHex := s.normalizeColorHex(req.ColorHex)
	normalizedMeaning := strings.TrimSpace(req.Meaning)

	for _, cm := range calendarColorMeanings {
		if cm.ID != colorMeaningID {
			if s.normalizeColorHex(cm.ColorHex) == normalizedColorHex {
				return nil, errors.New("color already exists for this calendar")
			}
			if strings.EqualFold(cm.Meaning, normalizedMeaning) {
				return nil, errors.New("meaning already exists for this calendar")
			}
		}
	}

	// Update color meaning
	updatedColorMeaning, err := s.queries.UpdateColorMeaning(ctx, db.UpdateColorMeaningParams{
		ID:       colorMeaningID,
		ColorHex: normalizedColorHex,
		Meaning:  normalizedMeaning,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update color meaning: %w", err)
	}

	return s.toColorMeaningResponse(updatedColorMeaning), nil
}

// DeleteColorMeaning deletes a color meaning
func (s *ColorMeaningService) DeleteColorMeaning(ctx context.Context, userID, colorMeaningID uuid.UUID) error {
	// Get color meaning and verify access
	_, err := s.GetColorMeaningByID(ctx, userID, colorMeaningID)
	if err != nil {
		return err
	}

	// Delete color meaning
	err = s.queries.DeleteColorMeaning(ctx, colorMeaningID)
	if err != nil {
		return fmt.Errorf("failed to delete color meaning: %w", err)
	}

	return nil
}

// Helper methods

func (s *ColorMeaningService) validateColorHex(colorHex string) error {
	if colorHex == "" {
		return ErrInvalidColorHex
	}

	// Remove # if present
	colorHex = strings.TrimPrefix(colorHex, "#")

	// Check if it's a valid hex color (3 or 6 characters)
	hexPattern := regexp.MustCompile(`^[0-9A-Fa-f]{3}$|^[0-9A-Fa-f]{6}$`)
	if !hexPattern.MatchString(colorHex) {
		return ErrInvalidColorHex
	}

	return nil
}

func (s *ColorMeaningService) validateMeaning(meaning string) error {
	meaning = strings.TrimSpace(meaning)
	if meaning == "" {
		return ErrMeaningEmpty
	}
	if len(meaning) > 50 {
		return errors.New("meaning cannot exceed 50 characters")
	}
	return nil
}

func (s *ColorMeaningService) normalizeColorHex(colorHex string) string {
	// Remove # if present and convert to uppercase
	colorHex = strings.TrimPrefix(colorHex, "#")
	colorHex = strings.ToUpper(colorHex)

	// Convert 3-character hex to 6-character hex
	if len(colorHex) == 3 {
		r := string(colorHex[0])
		g := string(colorHex[1])
		b := string(colorHex[2])
		colorHex = r + r + g + g + b + b
	}

	return "#" + colorHex
}

func (s *ColorMeaningService) toColorMeaningResponse(cm db.ColorMeaning) *ColorMeaningResponse {
	var createdAt string
	if cm.CreatedAt.Valid {
		createdAt = cm.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
	}

	return &ColorMeaningResponse{
		ID:         cm.ID,
		CalendarID: cm.CalendarID,
		ColorHex:   cm.ColorHex,
		Meaning:    cm.Meaning,
		CreatedAt:  createdAt,
	}
}
