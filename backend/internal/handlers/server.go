package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"days/internal/services"
)

type Server struct {
	userHandler     *UserHandler
	calendarHandler *CalendarHandler
}

func NewServer(
	userService services.UserServiceInterface,
	calendarService *services.CalendarService,
) *Server {
	return &Server{
		userHandler:     NewUserHandler(userService),
		calendarHandler: NewCalendarHandler(calendarService),
	}
}

func (s *Server) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Auth routes (no auth required) with body size limits (1MB)
	mux.HandleFunc("/api/users", CORSMiddleware(MaxBodyBytes(1<<20, s.userHandler.CreateUser)))
	mux.HandleFunc("/api/auth/login", CORSMiddleware(MaxBodyBytes(1<<20, s.userHandler.Login)))

	// Protected routes
	mux.HandleFunc("/api/users/", CORSMiddleware(AuthMiddleware(s.userHandler.GetUser)))
	mux.HandleFunc("/api/calendars", CORSMiddleware(AuthMiddleware(MaxBodyBytes(1<<20, s.handleCalendars))))
	mux.HandleFunc("/api/calendars/", CORSMiddleware(AuthMiddleware(MaxBodyBytes(1<<20, s.handleCalendarByID))))

	return mux
}

// handleCalendars routes requests to /api/calendars
func (s *Server) handleCalendars(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.calendarHandler.GetCalendars(w, r)
	case http.MethodPost:
		s.calendarHandler.CreateCalendar(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleCalendarByID routes requests to /api/calendars/{id}
func (s *Server) handleCalendarByID(w http.ResponseWriter, r *http.Request) {
	// Extract the path after /api/calendars/
	path := strings.TrimPrefix(r.URL.Path, "/api/calendars/")

	// If there's no ID, return 404
	if path == "" || path == "/" {
		http.Error(w, "Calendar ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.calendarHandler.GetCalendar(w, r)
	case http.MethodPut:
		s.calendarHandler.UpdateCalendar(w, r)
	case http.MethodDelete:
		s.calendarHandler.DeleteCalendar(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
