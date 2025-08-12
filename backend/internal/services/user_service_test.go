package services

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"days/internal/db"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockQueries implements a mock for the UserRepository interface
type MockQueries struct {
	mock.Mock
}

func (m *MockQueries) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockQueries) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockQueries) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.User), args.Error(1)
}

// Helper function to create a test user
func createTestUser(id uuid.UUID, email string) db.User {
	return db.User{
		ID:           id,
		Email:        email,
		PasswordHash: "hashed_password",
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}
}

func TestUserService_validateEmail(t *testing.T) {
	service := &UserService{}

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"valid email with subdomain", "user@mail.example.com", false},
		{"empty email", "", true},
		{"whitespace only", "   ", true},
		{"invalid format", "notanemail", true},
		{"missing @", "testexample.com", true},
		{"missing domain", "test@", true},
		{"multiple @", "test@@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidEmail, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_validatePassword(t *testing.T) {
	service := &UserService{}

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "password123", false},
		{"exactly 8 chars", "12345678", false},
		{"long password", "this_is_a_very_long_password_with_many_characters", false},
		{"too short", "1234567", true},
		{"empty password", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrWeakPassword, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_hashPassword(t *testing.T) {
	service := &UserService{}
	password := "testpassword123"

	hash1, err := service.hashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash1)
	assert.Contains(t, hash1, ":") // Should contain salt:hash format

	// Hash same password again - should be different due to random salt
	hash2, err := service.hashPassword(password)
	require.NoError(t, err)
	assert.NotEqual(t, hash1, hash2)

	// Both hashes should verify successfully
	assert.True(t, service.verifyPassword(password, hash1))
	assert.True(t, service.verifyPassword(password, hash2))
}

func TestUserService_verifyPassword(t *testing.T) {
	service := &UserService{}
	password := "correctpassword"

	hash, err := service.hashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{"correct password", password, hash, true},
		{"wrong password", "wrongpassword", hash, false},
		{"empty password", "", hash, false},
		{"malformed hash", password, "malformed", false},
		{"hash without colon", password, "nocolonhash", false},
		{"hash with invalid salt", password, "invalidsalt:validhash", false},
		{"hash with invalid hash part", password, "76616c6964:invalidhash", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.verifyPassword(tt.password, tt.hash)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()
	mockQueries := new(MockQueries)
	service := NewUserService(mockQueries)

	t.Run("successful user creation", func(t *testing.T) {
		req := CreateUserRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		userID := uuid.New()
		expectedUser := createTestUser(userID, "test@example.com")

		// Mock email check (user doesn't exist)
		mockQueries.On("GetUserByEmail", ctx, "test@example.com").
			Return(db.User{}, sql.ErrNoRows).Once()

		// Mock user creation
		mockQueries.On("CreateUser", ctx, mock.MatchedBy(func(params db.CreateUserParams) bool {
			return params.Email == "test@example.com" && params.PasswordHash != ""
		})).Return(expectedUser, nil).Once()

		result, err := service.CreateUser(ctx, req)

		require.NoError(t, err)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, "test@example.com", result.Email)
		assert.NotEmpty(t, result.CreatedAt)
		mockQueries.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		req := CreateUserRequest{
			Email:    "invalid-email",
			Password: "password123",
		}

		result, err := service.CreateUser(ctx, req)

		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidEmail, err)
		mockQueries.AssertExpectations(t)
	})

	t.Run("weak password", func(t *testing.T) {
		req := CreateUserRequest{
			Email:    "test@example.com",
			Password: "123", // Too short
		}

		result, err := service.CreateUser(ctx, req)

		assert.Nil(t, result)
		assert.Equal(t, ErrWeakPassword, err)
		mockQueries.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		req := CreateUserRequest{
			Email:    "existing@example.com",
			Password: "password123",
		}

		existingUser := createTestUser(uuid.New(), "existing@example.com")

		// Mock email check (user exists)
		mockQueries.On("GetUserByEmail", ctx, "existing@example.com").
			Return(existingUser, nil).Once()

		result, err := service.CreateUser(ctx, req)

		assert.Nil(t, result)
		assert.Equal(t, ErrEmailExists, err)
		mockQueries.AssertExpectations(t)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	ctx := context.Background()
	mockQueries := new(MockQueries)
	service := NewUserService(mockQueries)

	t.Run("user found", func(t *testing.T) {
		userID := uuid.New()
		expectedUser := createTestUser(userID, "test@example.com")

		mockQueries.On("GetUserByID", ctx, userID).
			Return(expectedUser, nil).Once()

		result, err := service.GetUserByID(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, "test@example.com", result.Email)
		mockQueries.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uuid.New()

		mockQueries.On("GetUserByID", ctx, userID).
			Return(db.User{}, sql.ErrNoRows).Once()

		result, err := service.GetUserByID(ctx, userID)

		assert.Nil(t, result)
		assert.Equal(t, ErrUserNotFound, err)
		mockQueries.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	ctx := context.Background()
	mockQueries := new(MockQueries)
	service := NewUserService(mockQueries)

	// Set up environment for JWT generation
	t.Setenv("JWT_SECRET", "test-secret-for-login")

	t.Run("successful login", func(t *testing.T) {
		userID := uuid.New()
		password := "correctpassword"
		email := "test@example.com"

		// Hash the password like the service would
		hash, err := service.hashPassword(password)
		require.NoError(t, err)

		user := db.User{
			ID:           userID,
			Email:        email,
			PasswordHash: hash,
			CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		}

		req := LoginRequest{
			Email:    email,
			Password: password,
		}

		mockQueries.On("GetUserByEmail", ctx, email).
			Return(user, nil).Once()

		result, err := service.Login(ctx, req)

		require.NoError(t, err)
		assert.Equal(t, userID, result.User.ID)
		assert.Equal(t, email, result.User.Email)
		assert.NotEmpty(t, result.Token)
		mockQueries.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		req := LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		mockQueries.On("GetUserByEmail", ctx, "nonexistent@example.com").
			Return(db.User{}, sql.ErrNoRows).Once()

		result, err := service.Login(ctx, req)

		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidCredentials, err)
		mockQueries.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		userID := uuid.New()
		email := "test@example.com"

		// Hash a different password
		hash, err := service.hashPassword("correctpassword")
		require.NoError(t, err)

		user := db.User{
			ID:           userID,
			Email:        email,
			PasswordHash: hash,
		}

		req := LoginRequest{
			Email:    email,
			Password: "wrongpassword",
		}

		mockQueries.On("GetUserByEmail", ctx, email).
			Return(user, nil).Once()

		result, err := service.Login(ctx, req)

		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidCredentials, err)
		mockQueries.AssertExpectations(t)
	})
}

func TestUserService_toUserResponse(t *testing.T) {
	service := &UserService{}
	userID := uuid.New()
	createdAt := time.Now()

	tests := []struct {
		name string
		user db.User
		want *UserResponse
	}{
		{
			name: "user with valid created_at",
			user: db.User{
				ID:        userID,
				Email:     "test@example.com",
				CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
			},
			want: &UserResponse{
				ID:        userID,
				Email:     "test@example.com",
				CreatedAt: createdAt.Format("2006-01-02T15:04:05Z"),
			},
		},
		{
			name: "user with invalid created_at",
			user: db.User{
				ID:        userID,
				Email:     "test@example.com",
				CreatedAt: sql.NullTime{Valid: false},
			},
			want: &UserResponse{
				ID:        userID,
				Email:     "test@example.com",
				CreatedAt: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.toUserResponse(tt.user)
			assert.Equal(t, tt.want, result)
		})
	}
}
