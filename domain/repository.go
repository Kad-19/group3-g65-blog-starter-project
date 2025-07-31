package domain

import (
	"context"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, Email string) (*User, error)
	UpdateUser(ctx context.Context, up UserProfile, Email string) error
	UpdateUserRole(ctx context.Context, role string, Email string) error
}
