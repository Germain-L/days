package dto

import "time"

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Username  string  `json:"username" binding:"required,min=3,max=50"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Timezone  *string `json:"timezone,omitempty"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Timezone  *string `json:"timezone,omitempty"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName *string   `json:"firstName,omitempty"`
	LastName  *string   `json:"lastName,omitempty"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UsersResponse represents a list of users with metadata
type UsersResponse struct {
	Data  []UserResponse `json:"data"`
	Total int            `json:"total"`
}
