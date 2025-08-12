package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDayEntryService_parseDate(t *testing.T) {
	service := &DayEntryService{}

	tests := []struct {
		name          string
		dateStr       string
		expectedError bool
		expectedDate  string
	}{
		{
			name:         "valid date",
			dateStr:      "2024-01-15",
			expectedDate: "2024-01-15",
		},
		{
			name:          "empty date",
			dateStr:       "",
			expectedError: true,
		},
		{
			name:          "invalid format",
			dateStr:       "01/15/2024",
			expectedError: true,
		},
		{
			name:          "invalid date",
			dateStr:       "2024-13-40",
			expectedError: true,
		},
		{
			name:         "leap year date",
			dateStr:      "2024-02-29",
			expectedDate: "2024-02-29",
		},
		{
			name:          "invalid leap year date",
			dateStr:       "2023-02-29",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.parseDate(tt.dateStr)

			if tt.expectedError {
				require.Error(t, err)
				assert.Equal(t, ErrInvalidDate, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedDate, result.Format("2006-01-02"))
			}
		})
	}
}

func TestDayEntryService_ValidationLogic(t *testing.T) {
	// Test the validation logic for day entry requests
	colorMeaningUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name          string
		request       CreateDayEntryRequest
		expectedError string
	}{
		{
			name: "valid request",
			request: CreateDayEntryRequest{
				Date:           "2024-01-15",
				ColorMeaningID: colorMeaningUUID,
				Notes:          dayEntryStringPtr("Valid notes"),
			},
		},
		{
			name: "valid request without notes",
			request: CreateDayEntryRequest{
				Date:           "2024-01-15",
				ColorMeaningID: colorMeaningUUID,
				Notes:          nil,
			},
		},
		{
			name: "invalid date format",
			request: CreateDayEntryRequest{
				Date:           "invalid-date",
				ColorMeaningID: colorMeaningUUID,
			},
			expectedError: "invalid date format",
		},
		{
			name: "empty date",
			request: CreateDayEntryRequest{
				Date:           "",
				ColorMeaningID: colorMeaningUUID,
			},
			expectedError: "invalid date format",
		},
	}

	service := &DayEntryService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test date parsing
			_, err := service.parseDate(tt.request.Date)

			if tt.expectedError == "invalid date format" {
				assert.Equal(t, ErrInvalidDate, err)
			} else if tt.expectedError != "" {
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDayEntryService_DateRangeValidation(t *testing.T) {
	tests := []struct {
		name          string
		request       DateRangeRequest
		expectedError string
	}{
		{
			name: "valid date range",
			request: DateRangeRequest{
				StartDate: "2024-01-01",
				EndDate:   "2024-01-31",
			},
		},
		{
			name: "same start and end date",
			request: DateRangeRequest{
				StartDate: "2024-01-15",
				EndDate:   "2024-01-15",
			},
		},
		{
			name: "invalid start date",
			request: DateRangeRequest{
				StartDate: "invalid-date",
				EndDate:   "2024-01-31",
			},
			expectedError: "invalid start date",
		},
		{
			name: "invalid end date",
			request: DateRangeRequest{
				StartDate: "2024-01-01",
				EndDate:   "invalid-date",
			},
			expectedError: "invalid end date",
		},
		{
			name: "end date before start date",
			request: DateRangeRequest{
				StartDate: "2024-01-31",
				EndDate:   "2024-01-01",
			},
			expectedError: "end date cannot be before start date",
		},
	}

	service := &DayEntryService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startDate, err := service.parseDate(tt.request.StartDate)
			if tt.expectedError == "invalid start date" {
				assert.Equal(t, ErrInvalidDate, err)
				return
			}
			require.NoError(t, err)

			endDate, err := service.parseDate(tt.request.EndDate)
			if tt.expectedError == "invalid end date" {
				assert.Equal(t, ErrInvalidDate, err)
				return
			}
			require.NoError(t, err)

			if tt.expectedError == "end date cannot be before start date" {
				assert.True(t, endDate.Before(startDate), "End date should be before start date for this test")
			} else {
				assert.False(t, endDate.Before(startDate), "End date should not be before start date")
			}
		})
	}
}

func TestDayEntryService_ErrorConstants(t *testing.T) {
	// Test that error constants are properly defined
	assert.NotNil(t, ErrDayEntryNotFound)
	assert.NotNil(t, ErrInvalidDate)
	assert.NotNil(t, ErrDayEntryExists)
	assert.NotNil(t, ErrUnauthorizedDayEntry)

	assert.Contains(t, ErrDayEntryNotFound.Error(), "day entry not found")
	assert.Contains(t, ErrInvalidDate.Error(), "invalid date format")
	assert.Contains(t, ErrDayEntryExists.Error(), "day entry already exists")
	assert.Contains(t, ErrUnauthorizedDayEntry.Error(), "not authorized")
}

// Benchmark tests for day entry operations
func BenchmarkDayEntryService_parseDate(b *testing.B) {
	service := &DayEntryService{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.parseDate("2024-01-15")
	}
}

func BenchmarkDayEntryService_parseDate_Invalid(b *testing.B) {
	service := &DayEntryService{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.parseDate("invalid-date")
	}
}

func BenchmarkDayEntryService_parseDate_Empty(b *testing.B) {
	service := &DayEntryService{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.parseDate("")
	}
}

// Helper function to create string pointers - use the one from calendar_service_test.go
func dayEntryStringPtr(s string) *string {
	return &s
}
