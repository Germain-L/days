package services

import (
	"context"
	"days/internal/db"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user database operations
type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
}

// CalendarRepository defines the interface for calendar database operations
type CalendarRepository interface {
	CreateCalendar(ctx context.Context, arg db.CreateCalendarParams) (db.Calendar, error)
	GetCalendarsByUserID(ctx context.Context, userID uuid.UUID) ([]db.Calendar, error)
	GetCalendarByID(ctx context.Context, id uuid.UUID) (db.Calendar, error)
	UpdateCalendar(ctx context.Context, arg db.UpdateCalendarParams) (db.Calendar, error)
	DeleteCalendar(ctx context.Context, id uuid.UUID) error
}

// UserServiceInterface defines the interface for user business logic
type UserServiceInterface interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}

// CalendarServiceInterface defines the interface for calendar business logic
type CalendarServiceInterface interface {
	CreateCalendar(ctx context.Context, userID uuid.UUID, req CreateCalendarRequest) (*CalendarResponse, error)
	GetCalendarsByUserID(ctx context.Context, userID uuid.UUID) ([]*CalendarResponse, error)
	GetCalendarByID(ctx context.Context, userID, calendarID uuid.UUID) (*CalendarResponse, error)
	UpdateCalendar(ctx context.Context, userID, calendarID uuid.UUID, req UpdateCalendarRequest) (*CalendarResponse, error)
	DeleteCalendar(ctx context.Context, userID, calendarID uuid.UUID) error
}

// Ensure db.Queries implements UserRepository
var _ UserRepository = (*db.Queries)(nil)

// Ensure UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)
