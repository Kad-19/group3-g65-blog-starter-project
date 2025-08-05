package domain

import (
	"context"
	"errors"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog) (string, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	UpdateBlog(ctx context.Context, blog *Blog, userid, id string) error
	DeleteBlog(ctx context.Context, id string) error
	ListBlogs(ctx context.Context, filter map[string]any, page, limit int) ([]*Blog, *Pagination, error)
}

type UserUseCase interface {
	Promote(ctx context.Context, email string) error
	Demote(ctx context.Context, email string) error
	ProfileUpdate(ctx context.Context, userid primitive.ObjectID, bio string, contactinfo string, file io.Reader) error
}

type AuthUsecase interface {
	Register(ctx context.Context, email, username, password string) (*User, error)
	Login(ctx context.Context, email, password string) (string, string, int, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, int, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID primitive.ObjectID) error
	ActivateUser(ctx context.Context, token string) error
	InitiateResetPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

var ErrUnauthorized = errors.New("unauthorized action")