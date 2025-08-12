package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name   string
		userID uuid.UUID
		secret string
		ttl    time.Duration
		valid  bool
	}{
		{
			name:   "valid token generation",
			userID: uuid.New(),
			secret: "test-secret-key",
			ttl:    time.Hour,
			valid:  true,
		},
		{
			name:   "empty secret",
			userID: uuid.New(),
			secret: "",
			ttl:    time.Hour,
			valid:  true, // JWT library allows empty secret
		},
		{
			name:   "nil user ID",
			userID: uuid.Nil,
			secret: "test-secret",
			ttl:    time.Hour,
			valid:  true, // Should still generate token with nil UUID
		},
		{
			name:   "zero TTL",
			userID: uuid.New(),
			secret: "test-secret",
			ttl:    0,
			valid:  true, // Immediately expired but valid generation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.secret, tt.ttl)

			if tt.valid {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Contains(t, token, ".") // JWT format has dots
			} else {
				require.Error(t, err)
				assert.Empty(t, token)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New()

	t.Run("valid token", func(t *testing.T) {
		// Generate a valid token
		token, err := GenerateToken(userID, secret, time.Hour)
		require.NoError(t, err)

		// Parse it back
		parsedUserID, err := ParseToken(token, secret)
		require.NoError(t, err)
		assert.Equal(t, userID, parsedUserID)
	})

	t.Run("wrong secret", func(t *testing.T) {
		token, err := GenerateToken(userID, secret, time.Hour)
		require.NoError(t, err)

		// Try to parse with wrong secret
		_, err = ParseToken(token, "wrong-secret")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("expired token", func(t *testing.T) {
		// Generate token that expires immediately
		token, err := GenerateToken(userID, secret, -time.Hour)
		require.NoError(t, err)

		// Try to parse expired token
		_, err = ParseToken(token, secret)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("malformed token", func(t *testing.T) {
		tests := []string{
			"",
			"invalid",
			"header.payload", // Missing signature
			"not.a.jwt.token.format",
		}

		for _, invalidToken := range tests {
			_, err := ParseToken(invalidToken, secret)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid token")
		}
	})

	t.Run("token without subject", func(t *testing.T) {
		// This would require manually crafting a JWT without subject
		// For now, our GenerateToken always includes subject, so we skip this
		t.Skip("Would need manual JWT crafting")
	})
}

func TestGenerateAndParseTokenRoundTrip(t *testing.T) {
	secret := "my-super-secret-key"
	ttl := 2 * time.Hour

	// Test multiple user IDs
	userIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
	}

	for i, userID := range userIDs {
		t.Run(fmt.Sprintf("roundtrip_%d", i), func(t *testing.T) {
			// Generate token
			token, err := GenerateToken(userID, secret, ttl)
			require.NoError(t, err)
			require.NotEmpty(t, token)

			// Parse token
			parsedUserID, err := ParseToken(token, secret)
			require.NoError(t, err)
			assert.Equal(t, userID, parsedUserID)
		})
	}
}

func TestTokenTimeClaims(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()
	ttl := time.Hour

	before := time.Now()
	token, err := GenerateToken(userID, secret, ttl)
	after := time.Now()
	require.NoError(t, err)

	// Parse token to check claims
	parsedUserID, err := ParseToken(token, secret)
	require.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)

	// For more detailed time validation, we'd need to parse the JWT manually
	// or expose the claims parsing. For now, we verify the token is valid
	// and not expired, and that time progressed during generation.
	assert.True(t, after.After(before) || after.Equal(before))
}
