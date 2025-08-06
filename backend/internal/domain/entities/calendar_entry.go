package entities

import (
	"days/internal/domain/valueobjects"
	"errors"
	"time"
)

// CalendarEntry represents a calendar entry entity
type CalendarEntry struct {
	id             *valueobjects.UserID // Reusing UserID type for UUIDs
	userID         *valueobjects.UserID
	date           *valueobjects.CalendarDate
	colorSettingID *valueobjects.UserID
	notes          *string
	createdAt      time.Time
	updatedAt      time.Time
}

// NewCalendarEntry creates a new CalendarEntry entity
func NewCalendarEntry(
	userID *valueobjects.UserID,
	date *valueobjects.CalendarDate,
	colorSettingID *valueobjects.UserID,
) (*CalendarEntry, error) {
	if userID == nil {
		return nil, errors.New("user ID is required")
	}
	if date == nil {
		return nil, errors.New("date is required")
	}
	if colorSettingID == nil {
		return nil, errors.New("color setting ID is required")
	}

	now := time.Now()
	return &CalendarEntry{
		id:             valueobjects.GenerateUserID(),
		userID:         userID,
		date:           date,
		colorSettingID: colorSettingID,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

// LoadCalendarEntry loads an existing calendar entry with all data
func LoadCalendarEntry(
	id *valueobjects.UserID,
	userID *valueobjects.UserID,
	date *valueobjects.CalendarDate,
	colorSettingID *valueobjects.UserID,
	notes *string,
	createdAt time.Time,
	updatedAt time.Time,
) *CalendarEntry {
	return &CalendarEntry{
		id:             id,
		userID:         userID,
		date:           date,
		colorSettingID: colorSettingID,
		notes:          notes,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// ID returns the calendar entry's ID
func (ce *CalendarEntry) ID() *valueobjects.UserID {
	return ce.id
}

// UserID returns the user ID that owns this calendar entry
func (ce *CalendarEntry) UserID() *valueobjects.UserID {
	return ce.userID
}

// Date returns the calendar entry's date
func (ce *CalendarEntry) Date() *valueobjects.CalendarDate {
	return ce.date
}

// ColorSettingID returns the color setting ID
func (ce *CalendarEntry) ColorSettingID() *valueobjects.UserID {
	return ce.colorSettingID
}

// Notes returns the calendar entry's notes
func (ce *CalendarEntry) Notes() *string {
	return ce.notes
}

// CreatedAt returns when the calendar entry was created
func (ce *CalendarEntry) CreatedAt() time.Time {
	return ce.createdAt
}

// UpdatedAt returns when the calendar entry was last updated
func (ce *CalendarEntry) UpdatedAt() time.Time {
	return ce.updatedAt
}

// UpdateColorSetting updates the calendar entry's color setting
func (ce *CalendarEntry) UpdateColorSetting(colorSettingID *valueobjects.UserID) error {
	if colorSettingID == nil {
		return errors.New("color setting ID is required")
	}
	ce.colorSettingID = colorSettingID
	ce.updatedAt = time.Now()
	return nil
}

// UpdateNotes updates the calendar entry's notes
func (ce *CalendarEntry) UpdateNotes(notes *string) {
	ce.notes = notes
	ce.updatedAt = time.Now()
}

// BelongsToUser checks if this calendar entry belongs to the specified user
func (ce *CalendarEntry) BelongsToUser(userID *valueobjects.UserID) bool {
	if userID == nil || ce.userID == nil {
		return false
	}
	return ce.userID.Equals(userID)
}

// IsForDate checks if this calendar entry is for the specified date
func (ce *CalendarEntry) IsForDate(date *valueobjects.CalendarDate) bool {
	if date == nil || ce.date == nil {
		return false
	}
	return ce.date.Equals(date)
}
