package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"days/internal/services"

	"github.com/google/uuid"
)

type CalendarHandler struct {
	calendarService *services.CalendarService
}

func NewCalendarHandler(calendarService *services.CalendarService) *CalendarHandler {
	return &CalendarHandler{
		calendarService: calendarService,
	}
}

// CreateCalendar handles POST /api/calendars
func (h *CalendarHandler) CreateCalendar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(ctxUserIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req services.CreateCalendarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	calendar, err := h.calendarService.CreateCalendar(r.Context(), userID, req)
	if err != nil {
		switch err {
		case services.ErrCalendarNameEmpty:
			writeJSONError(w, http.StatusBadRequest, err.Error())
		case services.ErrCalendarNameExists:
			writeJSONError(w, http.StatusConflict, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(calendar)
}

// GetCalendars handles GET /api/calendars
func (h *CalendarHandler) GetCalendars(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(ctxUserIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	calendars, err := h.calendarService.GetCalendarsByUserID(r.Context(), userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calendars)
}

// GetCalendar handles GET /api/calendars/{id}
func (h *CalendarHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(ctxUserIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Extract calendar ID from URL path
	calendarIDStr := extractIDFromPath(r.URL.Path, "/api/calendars/")
	calendarID, err := uuid.Parse(calendarIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid calendar ID")
		return
	}

	calendar, err := h.calendarService.GetCalendarByID(r.Context(), userID, calendarID)
	if err != nil {
		switch err {
		case services.ErrCalendarNotFound:
			writeJSONError(w, http.StatusNotFound, err.Error())
		case services.ErrUnauthorizedCalendar:
			writeJSONError(w, http.StatusForbidden, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calendar)
}

// UpdateCalendar handles PUT /api/calendars/{id}
func (h *CalendarHandler) UpdateCalendar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(ctxUserIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Extract calendar ID from URL path
	calendarIDStr := extractIDFromPath(r.URL.Path, "/api/calendars/")
	calendarID, err := uuid.Parse(calendarIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid calendar ID")
		return
	}

	var req services.UpdateCalendarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	calendar, err := h.calendarService.UpdateCalendar(r.Context(), userID, calendarID, req)
	if err != nil {
		switch err {
		case services.ErrCalendarNotFound:
			writeJSONError(w, http.StatusNotFound, err.Error())
		case services.ErrUnauthorizedCalendar:
			writeJSONError(w, http.StatusForbidden, err.Error())
		case services.ErrCalendarNameEmpty:
			writeJSONError(w, http.StatusBadRequest, err.Error())
		case services.ErrCalendarNameExists:
			writeJSONError(w, http.StatusConflict, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calendar)
}

// DeleteCalendar handles DELETE /api/calendars/{id}
func (h *CalendarHandler) DeleteCalendar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(ctxUserIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Extract calendar ID from URL path
	calendarIDStr := extractIDFromPath(r.URL.Path, "/api/calendars/")
	calendarID, err := uuid.Parse(calendarIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid calendar ID")
		return
	}

	err = h.calendarService.DeleteCalendar(r.Context(), userID, calendarID)
	if err != nil {
		switch err {
		case services.ErrCalendarNotFound:
			writeJSONError(w, http.StatusNotFound, err.Error())
		case services.ErrUnauthorizedCalendar:
			writeJSONError(w, http.StatusForbidden, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to extract ID from URL path
func extractIDFromPath(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}

	remaining := path[len(prefix):]

	// Handle potential sub-paths by taking only the first segment
	if idx := strings.Index(remaining, "/"); idx != -1 {
		return remaining[:idx]
	}

	return remaining
}
