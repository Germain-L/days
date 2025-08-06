package dto

import "time"

// CreateColorSettingRequest represents a request to create a color setting
type CreateColorSettingRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	HexColor    string  `json:"hexColor" binding:"required"`
	Description *string `json:"description,omitempty"`
	SortOrder   *int    `json:"sortOrder,omitempty"`
}

// UpdateColorSettingRequest represents a request to update a color setting
type UpdateColorSettingRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	HexColor    *string `json:"hexColor,omitempty"`
	Description *string `json:"description,omitempty"`
	SortOrder   *int    `json:"sortOrder,omitempty"`
}

// ColorSettingResponse represents a color setting in API responses
type ColorSettingResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	HexColor    string    `json:"hexColor"`
	Description *string   `json:"description,omitempty"`
	IsDefault   bool      `json:"isDefault"`
	SortOrder   int       `json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
