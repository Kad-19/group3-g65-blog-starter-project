package domain

import (
	"context"
)

type UserOperations interface {
	Promote(ctx context.Context, email string) error
	Demote(ctx context.Context, email string) error
	ProfileUpdate(ctx context.Context, up *UserProfile, email string) error
}
