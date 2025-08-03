package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog) (string, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	UpdateBlog(ctx context.Context, blog *Blog) error
	DeleteBlog(ctx context.Context, id string) error
	ListBlogs(ctx context.Context, filter map[string]interface{}) ([]*Blog, error)
}

type UserOperations interface {
	Promote(ctx context.Context, email string) error
	Demote(ctx context.Context, email string) error
	ProfileUpdate(ctx context.Context, up *UserProfile, email string) error
}

type AuthUsecase interface {
	Register(ctx context.Context, email, username, password string) (*User, error)
	Login(ctx context.Context, email, password string) (string, string, int, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, int, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID primitive.ObjectID) error
}