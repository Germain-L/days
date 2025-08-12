package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// AuthMiddleware is a simple middleware for demonstration
// In production, you'd implement proper JWT token validation
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check Bearer token format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "Token required", http.StatusUnauthorized)
			return
		}

		// For demonstration purposes, we'll use a simple token-to-userID mapping
		// In production, you'd validate JWT and extract user ID from claims
		userID := parseTokenToUserID(token)
		if userID == uuid.Nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// parseTokenToUserID is a placeholder for proper JWT validation
// In production, you'd decode the JWT and extract the user ID from claims
func parseTokenToUserID(token string) uuid.UUID {
	// For demonstration, we'll use a simple mapping
	// In production, this would be proper JWT validation
	switch token {
	case "test-token-user1":
		// Return a test user ID
		if id, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000"); err == nil {
			return id
		}
	case "test-token-user2":
		if id, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440001"); err == nil {
			return id
		}
	}

	// Try to parse the token as a UUID directly (for testing)
	if id, err := uuid.Parse(token); err == nil {
		return id
	}

	return uuid.Nil
}

// CORS middleware for development
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
