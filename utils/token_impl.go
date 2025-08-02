package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"g3-g65-bsp/domain"
	"time"
)

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func CreateActivationToken(email string, expiryDuration time.Duration) (*domain.ActivationToken, error) {
	// Generate a 32-byte long token (64 hex characters).
	tokenValue, err := generateSecureToken(32)
	if err != nil {
		return &domain.ActivationToken{}, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create the token object with an expiration time.
	// This is important for security to prevent tokens from being valid indefinitely.
	newToken := domain.ActivationToken{
		Email:     email,
		Token:     tokenValue,
		ExpiresAt: time.Now().Add(expiryDuration),
	}

	// Store the token in our "database".
	return &newToken, nil
}

func CreateResetToken(email string, expiryDuration time.Duration) (*domain.PasswordResetToken, error) {
	// Generate a 32-byte long token (64 hex characters).
	tokenValue, err := generateSecureToken(32)
	if err != nil {
		return &domain.PasswordResetToken{}, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create the token object with an expiration time.
	// This is important for security to prevent tokens from being valid indefinitely.
	newToken := domain.PasswordResetToken{
		Email:     email,
		Token:     tokenValue,
		ExpiresAt: time.Now().Add(expiryDuration),
	}

	// Store the token in our "database".
	return &newToken, nil
}
