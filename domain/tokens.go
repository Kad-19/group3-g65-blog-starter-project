package domain

import (
	"time"
)

type ActivationToken struct {
	Token     string
	Email     string
	ExpiresAt time.Time
}

type PasswordResetToken struct {
	Token     string
	Email     string
	ExpiresAt time.Time
}
