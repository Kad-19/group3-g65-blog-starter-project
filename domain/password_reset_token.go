package domain

import "time"

type PasswordResetToken struct {
	Email     string
	Token     string
	ExpiresAt time.Time
}
