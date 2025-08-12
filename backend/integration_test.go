package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"days/internal/auth"
	"days/internal/database"
	"days/internal/handlers"
	"days/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite defines our test suite
type IntegrationTestSuite struct {
	suite.Suite
	server     *handlers.Server
	httpServer *httptest.Server
	db         *database.Database
	userID     uuid.UUID
	token      string
}

// SetupSuite runs once before all tests
func (suite *IntegrationTestSuite) SetupSuite() {
	// Set up test environment
	os.Setenv("JWT_SECRET", "test-secret-for-integration")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "days_test")
	os.Setenv("DB_SSLMODE", "disable")

	// Connect to test database
	config := database.NewConfig()
	db, err := database.Connect(config)
	if err != nil {
		suite.T().Skip("Test database not available, skipping integration tests")
		return
	}
	suite.db = db

	// Initialize services
	userService := services.NewUserService(db.Queries)
	calendarService := services.NewCalendarService(db.Queries)

	// Initialize server
	suite.server = handlers.NewServer(userService, calendarService)
	mux := suite.server.SetupRoutes()
	suite.httpServer = httptest.NewServer(mux)
}

// TearDownSuite runs once after all tests
func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.httpServer != nil {
		suite.httpServer.Close()
	}
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest runs before each test
func (suite *IntegrationTestSuite) SetupTest() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
		return
	}

	// Clean up test data before each test
	suite.cleanupTestData()
}

// TearDownTest runs after each test
func (suite *IntegrationTestSuite) TearDownTest() {
	if suite.db == nil {
		return
	}

	// Clean up test data after each test
	suite.cleanupTestData()
}

func (suite *IntegrationTestSuite) cleanupTestData() {
	// Delete test users (cascades to calendars, etc.)
	_, err := suite.db.DB.Exec("DELETE FROM users WHERE email LIKE '%@test.integration'")
	if err != nil {
		suite.T().Logf("Warning: failed to clean up test data: %v", err)
	}
}

func (suite *IntegrationTestSuite) createTestUser() (uuid.UUID, string) {
	// Create a test user
	createReq := services.CreateUserRequest{
		Email:    fmt.Sprintf("user%d@test.integration", time.Now().UnixNano()),
		Password: "testpassword123",
	}

	reqBody, _ := json.Marshal(createReq)
	resp, err := http.Post(suite.httpServer.URL+"/api/users", "application/json", bytes.NewBuffer(reqBody))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	require.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var userResp services.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	require.NoError(suite.T(), err)

	// Login to get token
	loginReq := services.LoginRequest{
		Email:    createReq.Email,
		Password: createReq.Password,
	}

	reqBody, _ = json.Marshal(loginReq)
	resp, err = http.Post(suite.httpServer.URL+"/api/auth/login", "application/json", bytes.NewBuffer(reqBody))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	require.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var loginResp services.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	require.NoError(suite.T(), err)

	return userResp.ID, loginResp.Token
}

func (suite *IntegrationTestSuite) TestHealthEndpoint() {
	resp, err := http.Get(suite.httpServer.URL + "/health")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	assert.Equal(suite.T(), "OK", string(body[:n]))
}

func (suite *IntegrationTestSuite) TestUserRegistrationAndLogin() {
	// Test user creation
	createReq := services.CreateUserRequest{
		Email:    "testuser@test.integration",
		Password: "testpassword123",
	}

	reqBody, _ := json.Marshal(createReq)
	resp, err := http.Post(suite.httpServer.URL+"/api/users", "application/json", bytes.NewBuffer(reqBody))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var userResp services.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), createReq.Email, userResp.Email)
	assert.NotEmpty(suite.T(), userResp.ID)

	// Test duplicate email
	reqBody, _ = json.Marshal(createReq)
	resp, err = http.Post(suite.httpServer.URL+"/api/users", "application/json", bytes.NewBuffer(reqBody))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusConflict, resp.StatusCode)

	// Test login with correct credentials
	loginReq := services.LoginRequest{
		Email:    createReq.Email,
		Password: createReq.Password,
	}

	reqBody, _ = json.Marshal(loginReq)
	resp, err = http.Post(suite.httpServer.URL+"/api/auth/login", "application/json", bytes.NewBuffer(reqBody))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var loginResp services.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), userResp.ID, loginResp.User.ID)
	assert.NotEmpty(suite.T(), loginResp.Token)

	// Verify token can be parsed
	secret := os.Getenv("JWT_SECRET")
	parsedUserID, err := auth.ParseToken(loginResp.Token, secret)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), userResp.ID, parsedUserID)

	// Test login with wrong password
	loginReq.Password = "wrongpassword"
	reqBody, _ = json.Marshal(loginReq)
	resp, err = http.Post(suite.httpServer.URL+"/api/auth/login", "application/json", bytes.NewBuffer(reqBody))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestAuthenticatedUserAccess() {
	userID, token := suite.createTestUser()

	// Test accessing own user record
	req, _ := http.NewRequest(http.MethodGet, suite.httpServer.URL+"/api/users/"+userID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var userResp services.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, userResp.ID)

	// Test accessing other user record (should fail)
	otherUserID := uuid.New()
	req, _ = http.NewRequest(http.MethodGet, suite.httpServer.URL+"/api/users/"+otherUserID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusForbidden, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestUnauthorizedAccess() {
	userID := uuid.New()

	tests := []struct {
		name   string
		url    string
		method string
		token  string
		status int
	}{
		{
			name:   "no token",
			url:    "/api/users/" + userID.String(),
			method: http.MethodGet,
			token:  "",
			status: http.StatusUnauthorized,
		},
		{
			name:   "invalid token",
			url:    "/api/users/" + userID.String(),
			method: http.MethodGet,
			token:  "invalid.token.here",
			status: http.StatusUnauthorized,
		},
		{
			name:   "malformed auth header",
			url:    "/api/users/" + userID.String(),
			method: http.MethodGet,
			token:  "NotBearer invalid.token.here",
			status: http.StatusUnauthorized,
		},
	}

	client := &http.Client{}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, suite.httpServer.URL+tt.url, nil)
			if tt.token != "" {
				if tt.token == "NotBearer invalid.token.here" {
					req.Header.Set("Authorization", tt.token)
				} else {
					req.Header.Set("Authorization", "Bearer "+tt.token)
				}
			}

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.status, resp.StatusCode)
		})
	}
}

func (suite *IntegrationTestSuite) TestBodySizeLimits() {
	// Create oversized request body (> 1MB)
	largeBody := bytes.Repeat([]byte("a"), 2*1024*1024) // 2MB

	req, _ := http.NewRequest(http.MethodPost, suite.httpServer.URL+"/api/users", bytes.NewReader(largeBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Should get a 400 or connection error due to body size limit
	assert.NotEqual(suite.T(), http.StatusCreated, resp.StatusCode)
}

// TestIntegrationSuite runs the integration test suite
func TestIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	suite.Run(t, new(IntegrationTestSuite))
}

// TestMain allows us to set up and tear down for all tests
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
}
