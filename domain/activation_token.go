package domain

import (
	"time"
)

type ActivationToken struct {
	Token     string
	Email     string
	ExpiresAt time.Time
}
