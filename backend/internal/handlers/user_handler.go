package handlers

import (
	"encoding/json"
	"net/http"

	"days/internal/services"

	"github.com/google/uuid"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser handles POST /api/users
//
//	@Summary		Create a new user
//	@Description	Create a new user account with email and password
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		services.CreateUserRequest	true	"User creation request"
//	@Success		201		{object}	services.UserResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req services.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		switch err {
		case services.ErrInvalidEmail:
			writeJSONError(w, http.StatusBadRequest, err.Error())
		case services.ErrWeakPassword:
			writeJSONError(w, http.StatusBadRequest, err.Error())
		case services.ErrEmailExists:
			writeJSONError(w, http.StatusConflict, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login handles POST /api/auth/login
//
//	@Summary		User login
//	@Description	Authenticate user and return JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		services.LoginRequest	true	"Login credentials"
//	@Success		200			{object}	services.LoginResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req services.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	loginResponse, err := h.userService.Login(r.Context(), req)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			writeJSONError(w, http.StatusUnauthorized, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

// GetUser handles GET /api/users/{id}
//
//	@Summary		Get user by ID
//	@Description	Retrieve user information by user ID (self only)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	services.UserResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from URL path (you'd typically use a router for this)
	userIDStr := r.URL.Path[len("/api/users/"):]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Enforce self-only access for now
	authUserID, ok := r.Context().Value(ctxUserIDKey).(uuid.UUID)
	if !ok || authUserID != userID {
		writeJSONError(w, http.StatusForbidden, "forbidden: can only access own user record")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
