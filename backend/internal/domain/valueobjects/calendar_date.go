package valueobjects

import (
	"errors"
	"time"
)

// CalendarDate represents a specific date in the calendar
type CalendarDate struct {
	value time.Time
}

// NewCalendarDate creates a new CalendarDate from a time.Time
func NewCalendarDate(date time.Time) (*CalendarDate, error) {
	if date.IsZero() {
		return nil, errors.New("date cannot be zero")
	}

	// Normalize to date only (remove time component)
	normalized := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	return &CalendarDate{value: normalized}, nil
}

// NewCalendarDateFromString creates a new CalendarDate from a string in YYYY-MM-DD format
func NewCalendarDateFromString(dateStr string) (*CalendarDate, error) {
	if dateStr == "" {
		return nil, errors.New("date string cannot be empty")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	return NewCalendarDate(date)
}

// Today creates a CalendarDate for today
func Today() *CalendarDate {
	now := time.Now()
	date, _ := NewCalendarDate(now) // Today() can't fail
	return date
}

// Value returns the time.Time value of the date
func (cd *CalendarDate) Value() time.Time {
	return cd.value
}

// String returns the date in YYYY-MM-DD format
func (cd *CalendarDate) String() string {
	return cd.value.Format("2006-01-02")
}

// Equals checks if two dates are equal
func (cd *CalendarDate) Equals(other *CalendarDate) bool {
	if other == nil {
		return false
	}
	return cd.value.Equal(other.value)
}

// Before checks if this date is before another date
func (cd *CalendarDate) Before(other *CalendarDate) bool {
	if other == nil {
		return false
	}
	return cd.value.Before(other.value)
}

// After checks if this date is after another date
func (cd *CalendarDate) After(other *CalendarDate) bool {
	if other == nil {
		return false
	}
	return cd.value.After(other.value)
}

// AddDays adds the specified number of days to the date
func (cd *CalendarDate) AddDays(days int) *CalendarDate {
	newDate := cd.value.AddDate(0, 0, days)
	result, _ := NewCalendarDate(newDate) // AddDays can't fail
	return result
}
