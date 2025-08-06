package dto

import "time"

// CreateCalendarEntryRequest represents a request to create a calendar entry
type CreateCalendarEntryRequest struct {
	Date           string  `json:"date" binding:"required"`
	ColorSettingID string  `json:"colorSettingId" binding:"required"`
	Notes          *string `json:"notes,omitempty"`
}

// UpdateCalendarEntryRequest represents a request to update a calendar entry
type UpdateCalendarEntryRequest struct {
	ColorSettingID *string `json:"colorSettingId,omitempty"`
	Notes          *string `json:"notes,omitempty"`
}

// CalendarEntryResponse represents a calendar entry in API responses
type CalendarEntryResponse struct {
	ID           string                `json:"id"`
	UserID       string                `json:"userId"`
	Date         string                `json:"date"`
	ColorSetting *ColorSettingResponse `json:"colorSetting,omitempty"`
	Notes        *string               `json:"notes,omitempty"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}

// CalendarEntriesResponse represents a list of calendar entries with metadata
type CalendarEntriesResponse struct {
	Data    []CalendarEntryResponse `json:"data"`
	Total   int                     `json:"total"`
	Filters CalendarEntryFilters    `json:"filters,omitempty"`
}

// CalendarEntryFilters represents the filters applied to calendar entry queries
type CalendarEntryFilters struct {
	StartDate      *string `json:"startDate,omitempty"`
	EndDate        *string `json:"endDate,omitempty"`
	ColorSettingID *string `json:"colorSettingId,omitempty"`
}
