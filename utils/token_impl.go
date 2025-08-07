package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"g3-g65-bsp/domain"
	"time"
)


func GenerateRandomToken() (string, *time.Time, error) {
	b := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(b); err != nil {
		return "", nil, errors.New("failed to generate random token")
	}
	token := base64.URLEncoding.EncodeToString(b)
	expiry := time.Now().Add(30 * time.Minute)
	return token, &expiry, nil
}

func CreateResetToken(email string, expiryDuration time.Duration) (*domain.PasswordResetToken, error) {
	tokenValue, expiry, err := GenerateRandomToken()
	if err != nil {
		return &domain.PasswordResetToken{}, fmt.Errorf("failed to generate token: %w", err)
	}

	newToken := domain.PasswordResetToken{
		Email:     email,
		Token:     tokenValue,
		ExpiresAt: *expiry,
	}
	return &newToken, nil
}
