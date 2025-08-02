package domain

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog) (string, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	UpdateBlog(ctx context.Context, blog *Blog) error
	DeleteBlog(ctx context.Context, id string) error
	ListBlogs(ctx context.Context, filter map[string]interface{}) ([]*Blog, error)
}

type UserUseCase interface {
	Promote(ctx context.Context, email string) error
	Demote(ctx context.Context, email string) error
	ProfileUpdate(ctx context.Context, userid primitive.ObjectID, bio string, contactinfo string, file io.Reader) error
}

type AuthUsecase interface {
	ActivateUser(ctx context.Context, token string) error
	Register(ctx context.Context, req UserCreateRequest) (*UserResponse, error)
	Login(ctx context.Context, email, password string) (string, string, int, error)
	InitiateResetPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}
