package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"days/internal/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteJSONError(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		message    string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "bad request error",
			status:     http.StatusBadRequest,
			message:    "invalid input",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid input"}`,
		},
		{
			name:       "unauthorized error",
			status:     http.StatusUnauthorized,
			message:    "token required",
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"error":"token required"}`,
		},
		{
			name:       "internal server error",
			status:     http.StatusInternalServerError,
			message:    "something went wrong",
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"something went wrong"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeJSONError(w, tt.status, tt.message)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			// Parse and compare JSON to avoid whitespace issues
			var got, want errorResponse
			err := json.Unmarshal(w.Body.Bytes(), &got)
			require.NoError(t, err)
			err = json.Unmarshal([]byte(tt.wantBody), &want)
			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	// Set up test environment
	secret := "test-secret-for-middleware"
	userID := uuid.New()

	// Generate a valid token
	validToken, err := auth.GenerateToken(userID, secret, time.Hour)
	require.NoError(t, err)

	// Create a test handler that checks the context
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		ctxUserID, ok := r.Context().Value(ctxUserIDKey).(uuid.UUID)
		if !ok {
			t.Error("userID not found in context")
			return
		}
		if ctxUserID != userID {
			t.Errorf("expected userID %v, got %v", userID, ctxUserID)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	tests := []struct {
		name           string
		jwtSecret      string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid token",
			jwtSecret:      secret,
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "missing authorization header",
			jwtSecret:      secret,
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"authorization header required"}`,
		},
		{
			name:           "invalid authorization format",
			jwtSecret:      secret,
			authHeader:     "InvalidFormat " + validToken,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid authorization format"}`,
		},
		{
			name:           "missing bearer token",
			jwtSecret:      secret,
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"token required"}`,
		},
		{
			name:           "invalid token",
			jwtSecret:      secret,
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid or expired token"}`,
		},
		{
			name:           "wrong secret",
			jwtSecret:      "wrong-secret",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid or expired token"}`,
		},
		{
			name:           "missing JWT secret in env",
			jwtSecret:      "",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"server misconfigured: missing JWT secret"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.jwtSecret != "" {
				t.Setenv("JWT_SECRET", tt.jwtSecret)
			} else {
				os.Unsetenv("JWT_SECRET")
			}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			handler := AuthMiddleware(testHandler)
			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if strings.Contains(tt.expectedBody, `{"error":`) {
				// Compare JSON responses
				var got, want errorResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				err = json.Unmarshal([]byte(tt.expectedBody), &want)
				require.NoError(t, err)
				assert.Equal(t, want, got)
			} else {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestMaxBodyBytes(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		// Try to read the body
		body := make([]byte, 1024)
		n, err := r.Body.Read(body)
		if err != nil && err.Error() != "EOF" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error reading body: " + err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("read %d bytes", n)))
	}

	tests := []struct {
		name         string
		limit        int64
		bodySize     int
		expectError  bool
		expectedCode int
	}{
		{
			name:         "body within limit",
			limit:        100,
			bodySize:     50,
			expectError:  false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "body at exact limit",
			limit:        100,
			bodySize:     100,
			expectError:  false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "body exceeds limit",
			limit:        100,
			bodySize:     150,
			expectError:  true,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.Repeat([]byte("a"), tt.bodySize)
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler := MaxBodyBytes(tt.limit, testHandler)
			handler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectError {
				assert.Contains(t, w.Body.String(), "error reading body")
			} else {
				assert.Contains(t, w.Body.String(), fmt.Sprintf("read %d bytes", tt.bodySize))
			}
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	tests := []struct {
		name            string
		corsAllowOrigin string
		method          string
		expectedStatus  int
		expectedOrigin  string
	}{
		{
			name:            "GET request with default CORS",
			corsAllowOrigin: "",
			method:          http.MethodGet,
			expectedStatus:  http.StatusOK,
			expectedOrigin:  "*",
		},
		{
			name:            "OPTIONS request",
			corsAllowOrigin: "",
			method:          http.MethodOptions,
			expectedStatus:  http.StatusOK,
			expectedOrigin:  "*",
		},
		{
			name:            "custom CORS origin",
			corsAllowOrigin: "https://example.com",
			method:          http.MethodGet,
			expectedStatus:  http.StatusOK,
			expectedOrigin:  "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.corsAllowOrigin != "" {
				t.Setenv("CORS_ALLOW_ORIGIN", tt.corsAllowOrigin)
			} else {
				os.Unsetenv("CORS_ALLOW_ORIGIN")
			}

			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			handler := CORSMiddleware(testHandler)
			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))

			if tt.method == http.MethodOptions {
				assert.Empty(t, w.Body.String())
			} else {
				assert.Equal(t, "success", w.Body.String())
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	tests := []struct {
		name          string
		timeout       time.Duration
		handlerDelay  time.Duration
		expectTimeout bool
	}{
		{
			name:          "request completes within timeout",
			timeout:       100 * time.Millisecond,
			handlerDelay:  50 * time.Millisecond,
			expectTimeout: false,
		},
		{
			name:          "request exceeds timeout",
			timeout:       50 * time.Millisecond,
			handlerDelay:  100 * time.Millisecond,
			expectTimeout: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler := func(w http.ResponseWriter, r *http.Request) {
				select {
				case <-time.After(tt.handlerDelay):
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("success"))
				case <-r.Context().Done():
					// Context was cancelled due to timeout
					return
				}
			}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			handler := WithTimeout(tt.timeout, testHandler)
			handler(w, req)

			if tt.expectTimeout {
				// The handler should not have written anything due to context cancellation
				assert.Equal(t, http.StatusOK, w.Code) // Default status
				assert.Empty(t, w.Body.String())
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, "success", w.Body.String())
			}
		})
	}
}
