package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHasher(t *testing.T) {
	hasher := &PasswordHasher{}
	password := "password123"

	// Test hashing
	hash, err := hasher.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Test correct password comparison
	assert.True(t, hasher.CompareHashAndPassword(hash, password))

	// Test incorrect password comparison
	assert.False(t, hasher.CompareHashAndPassword(hash, "wrongpassword"))
}
