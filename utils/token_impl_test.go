package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomToken(t *testing.T) {
	token, expiry, err := GenerateRandomToken()

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotNil(t, expiry)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), *expiry, time.Second)
}

func TestCreateResetToken(t *testing.T) {
	email := "test@example.com"
	resetToken, err := CreateResetToken(email, time.Hour)

	assert.NoError(t, err)
	assert.NotNil(t, resetToken)
	assert.Equal(t, email, resetToken.Email)
	assert.NotEmpty(t, resetToken.Token)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), resetToken.ExpiresAt, time.Second)
}
