package handlers

import (
	"encoding/json"
	"net/http"

	"days/internal/services"

	"github.com/google/uuid"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser handles POST /api/users
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
