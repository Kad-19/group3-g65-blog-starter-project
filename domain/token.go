package domain

import "time"

type PasswordResetToken struct {
	Email     string
	Token     string
	ExpiresAt time.Time
}


type RefreshToken struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}