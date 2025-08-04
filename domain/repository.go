package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogRepository interface {
	CreateBlog(ctx context.Context, blog *Blog) (string, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	UpdateBlog(ctx context.Context, blog *Blog) error
	DeleteBlog(ctx context.Context, id string) error
	ListBlogs(ctx context.Context, filter map[string]interface{}) ([]*Blog, error)
	IncrementBlogViewCount(ctx context.Context, id string, blog *Blog) error
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	UpdateUser(ctx context.Context, bio string, contactInfo string, imagePath string, Email string) error
	UpdateUserRole(ctx context.Context, role string, Email string) error
	UpdateActiveStatus(ctx context.Context, email string) error
	UpdateUserPassword(ctx context.Context, email string, newPasswordHash string) error
}

type ActivationTokenRepository interface {
	Create(ctx context.Context, token *ActivationToken) error
	GetByToken(ctx context.Context, token string) (*ActivationToken, error)
	Delete(ctx context.Context, token string) error
}

type PasswordResetRepository interface {
	Create(ctx context.Context, token *PasswordResetToken) error
	GetByToken(ctx context.Context, token string) (*PasswordResetToken, error)
	Delete(ctx context.Context, token string) error
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID primitive.ObjectID, tokenHash string, expiresAt time.Time) error
	FindAndDeleteRefreshToken(ctx context.Context, tokenHash string) (primitive.ObjectID, error)
	DeleteAllForUser(ctx context.Context, userID primitive.ObjectID) error
}
