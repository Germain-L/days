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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req services.CreateCalendarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	calendar, err := h.calendarService.CreateCalendar(r.Context(), userID, req)
	if err != nil {
		switch err {
		case services.ErrCalendarNameEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case services.ErrCalendarNameExists:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	calendars, err := h.calendarService.GetCalendarsByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calendars)
}

// GetCalendar handles GET /api/calendars/{id}
func (h *CalendarHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract calendar ID from URL path
	calendarIDStr := extractIDFromPath(r.URL.Path, "/api/calendars/")
	calendarID, err := uuid.Parse(calendarIDStr)
	if err != nil {
		http.Error(w, "Invalid calendar ID", http.StatusBadRequest)
		return
	}

	calendar, err := h.calendarService.GetCalendarByID(r.Context(), userID, calendarID)
	if err != nil {
		switch err {
		case services.ErrCalendarNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case services.ErrUnauthorizedCalendar:
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calendar)
}

// UpdateCalendar handles PUT /api/calendars/{id}
func (h *CalendarHandler) UpdateCalendar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract calendar ID from URL path
	calendarIDStr := extractIDFromPath(r.URL.Path, "/api/calendars/")
	calendarID, err := uuid.Parse(calendarIDStr)
	if err != nil {
		http.Error(w, "Invalid calendar ID", http.StatusBadRequest)
		return
	}

	var req services.UpdateCalendarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	calendar, err := h.calendarService.UpdateCalendar(r.Context(), userID, calendarID, req)
	if err != nil {
		switch err {
		case services.ErrCalendarNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case services.ErrUnauthorizedCalendar:
			http.Error(w, err.Error(), http.StatusForbidden)
		case services.ErrCalendarNameEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case services.ErrCalendarNameExists:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calendar)
}

// DeleteCalendar handles DELETE /api/calendars/{id}
func (h *CalendarHandler) DeleteCalendar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract calendar ID from URL path
	calendarIDStr := extractIDFromPath(r.URL.Path, "/api/calendars/")
	calendarID, err := uuid.Parse(calendarIDStr)
	if err != nil {
		http.Error(w, "Invalid calendar ID", http.StatusBadRequest)
		return
	}

	err = h.calendarService.DeleteCalendar(r.Context(), userID, calendarID)
	if err != nil {
		switch err {
		case services.ErrCalendarNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case services.ErrUnauthorizedCalendar:
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
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
