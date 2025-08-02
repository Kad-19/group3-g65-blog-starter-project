package domain

import "time"

type PasswordResetToken struct {
	Token     string
	Email     string
	ExpiresAt time.Time
}
