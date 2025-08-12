package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"days/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserService implements a mock for the UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req services.CreateUserRequest) (*services.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*services.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.UserResponse), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, req services.LoginRequest) (*services.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.LoginResponse), args.Error(1)
}

func TestUserHandler_CreateUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("successful user creation", func(t *testing.T) {
		req := services.CreateUserRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		userID := uuid.New()
		expectedResponse := &services.UserResponse{
			ID:        userID,
			Email:     "test@example.com",
			CreatedAt: "2023-01-01T00:00:00Z",
		}

		mockService.On("CreateUser", mock.Anything, req).Return(expectedResponse, nil).Once()

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateUser(w, httpReq)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response services.UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Email, response.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid method", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		w := httptest.NewRecorder()

		handler.CreateUser(w, httpReq)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp errorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		require.NoError(t, err)
		assert.Equal(t, "method not allowed", errorResp.Error)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader("invalid json"))
		w := httptest.NewRecorder()

		handler.CreateUser(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResp errorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		require.NoError(t, err)
		assert.Equal(t, "invalid JSON", errorResp.Error)
	})

	t.Run("invalid email", func(t *testing.T) {
		req := services.CreateUserRequest{
			Email:    "invalid-email",
			Password: "password123",
		}

		mockService.On("CreateUser", mock.Anything, req).Return(nil, services.ErrInvalidEmail).Once()

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler.CreateUser(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResp errorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		require.NoError(t, err)
		assert.Equal(t, services.ErrInvalidEmail.Error(), errorResp.Error)

		mockService.AssertExpectations(t)
	})

	t.Run("weak password", func(t *testing.T) {
		req := services.CreateUserRequest{
			Email:    "test@example.com",
			Password: "123",
		}

		mockService.On("CreateUser", mock.Anything, req).Return(nil, services.ErrWeakPassword).Once()

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler.CreateUser(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResp errorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		require.NoError(t, err)
		assert.Equal(t, services.ErrWeakPassword.Error(), errorResp.Error)

		mockService.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		req := services.CreateUserRequest{
			Email:    "existing@example.com",
			Password: "password123",
		}

		mockService.On("CreateUser", mock.Anything, req).Return(nil, services.ErrEmailExists).Once()

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler.CreateUser(w, httpReq)

		assert.Equal(t, http.StatusConflict, w.Code)

		var errorResp errorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		require.NoError(t, err)
		assert.Equal(t, services.ErrEmailExists.Error(), errorResp.Error)

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Login(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("successful login", func(t *testing.T) {
		req := services.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		userID := uuid.New()
		expectedResponse := &services.LoginResponse{
			User: services.UserResponse{
				ID:        userID,
				Email:     "test@example.com",
				CreatedAt: "2023-01-01T00:00:00Z",
			},
			Token: "jwt.token.here",
		}

		mockService.On("Login", mock.Anything, req).Return(expectedResponse, nil).Once()

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler.Login(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response services.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, expectedResponse.User.ID, response.User.ID)
		assert.Equal(t, expectedResponse.Token, response.Token)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		req := services.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		mockService.On("Login", mock.Anything, req).Return(nil, services.ErrInvalidCredentials).Once()

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler.Login(w, httpReq)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var errorResp errorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		require.NoError(t, err)
		assert.Equal(t, services.ErrInvalidCredentials.Error(), errorResp.Error)

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	t.Run("successful get user (self)", func(t *testing.T) {
		expectedResponse := &services.UserResponse{
			ID:        userID,
			Email:     "test@example.com",
			CreatedAt: "2023-01-01T00:00:00Z",
		}

		mockService.On("GetUserByID", mock.Anything, userID).Return(expectedResponse, nil).Once()

		httpReq := httptest.NewRequest(http.MethodGet, "/api/users/"+userID.String(), nil)
		// Set the authenticated user ID in context
		ctx := context.WithValue(httpReq.Context(), ctxUserIDKey, userID)
		httpReq = httpReq.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetUser(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response services.UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Email, response.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("forbidden - trying to access other user", func(t *testing.T) {
		otherUserID := uuid.New()
		httpReq := httptest.NewRequest(http.MethodGet, "/api/users/"+otherUserID.String(), nil)
		// Set different authenticated user ID in context
		ctx := context.WithValue(httpReq.Context(), ctxUserIDKey, userID)
		httpReq = httpReq.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetUser(w, httpReq)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var errorResp errorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		require.NoError(t, err)
		assert.Equal(t, "forbidden: can only access own user record", errorResp.Error)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodGet, "/api/users/invalid-uuid", nil)
		ctx := context.WithValue(httpReq.Context(), ctxUserIDKey, userID)
		httpReq = httpReq.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetUser(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResp errorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		require.NoError(t, err)
		assert.Equal(t, "invalid user ID", errorResp.Error)
	})

	t.Run("user not found", func(t *testing.T) {
		mockService.On("GetUserByID", mock.Anything, userID).Return(nil, services.ErrUserNotFound).Once()

		httpReq := httptest.NewRequest(http.MethodGet, "/api/users/"+userID.String(), nil)
		ctx := context.WithValue(httpReq.Context(), ctxUserIDKey, userID)
		httpReq = httpReq.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetUser(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var errorResp errorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		require.NoError(t, err)
		assert.Equal(t, services.ErrUserNotFound.Error(), errorResp.Error)

		mockService.AssertExpectations(t)
	})

	t.Run("missing user ID in context", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodGet, "/api/users/"+userID.String(), nil)
		// No user ID in context (simulates missing auth)
		w := httptest.NewRecorder()

		handler.GetUser(w, httpReq)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var errorResp errorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		require.NoError(t, err)
		assert.Equal(t, "forbidden: can only access own user record", errorResp.Error)
	})
}
