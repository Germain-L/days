package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"days/internal/auth"

	"github.com/google/uuid"
)

// typed context key to avoid collisions
type ctxKey string

const ctxUserIDKey ctxKey = "userID"

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: msg})
}

// AuthMiddleware validates JWT Bearer tokens and injects user ID into context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, http.StatusUnauthorized, "authorization header required")
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeJSONError(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if token == "" {
			writeJSONError(w, http.StatusUnauthorized, "token required")
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			writeJSONError(w, http.StatusInternalServerError, "server misconfigured: missing JWT secret")
			return
		}

		userID, err := auth.ParseToken(token, secret)
		if err != nil || userID == uuid.Nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// MaxBodyBytes limits the size of request bodies to prevent DoS via large payloads
func MaxBodyBytes(n int64, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, n)
		next.ServeHTTP(w, r)
	}
}

// CORS middleware for development (restrict in production via env)
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")
		if allowOrigin == "" {
			allowOrigin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// A small timeout wrapper to help ensure handlers don't hang indefinitely
func WithTimeout(d time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), d)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
