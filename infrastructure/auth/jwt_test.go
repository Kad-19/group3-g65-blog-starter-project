package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	jwt := NewJWT("access-secret", "refresh-secret", time.Minute*15, time.Hour*24)
	userID := "user123"
	role := "user"

	// Test access token generation
	accessToken, err := jwt.GenerateAccessToken(userID, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)

	// Test access token validation
	claims, err := jwt.ValidateAccessToken(accessToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)

	// Test refresh token generation
	refreshToken, err := jwt.GenerateRefreshToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)
}
