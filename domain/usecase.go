package domain

import (
	"context"
	"errors"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog, userid string) (*Blog, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	UpdateBlog(ctx context.Context, blog *Blog, userid, id string) (*Blog, error)
	DeleteBlog(ctx context.Context, id string) error
	ListBlogs(ctx context.Context, filter map[string]any, page, limit int) ([]*Blog, *Pagination, error)
}

type UserUsecase interface {
	Promote(ctx context.Context, email string) error
	Demote(ctx context.Context, email string) error
	ProfileUpdate(ctx context.Context, userid primitive.ObjectID, bio string, contactinfo string, file io.Reader) error
	GetAllUsers(ctx context.Context, page int, limit int) ([]User, Pagination, error)
}

type AuthUsecase interface {
	Register(ctx context.Context, email, username, password string) error
	Login(ctx context.Context, email, password string) (string, string, int, *User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, int, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID primitive.ObjectID) error
	ActivateUser(ctx context.Context, token, email string) error
	ResendActivationEmail(ctx context.Context, email string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type InteractionUsecase interface {
	LikeBlog(ctx context.Context, userID string, blogID string, preftype string) error
	CommentOnBlog(ctx context.Context, userID string, blogID string, comment *Comment) error
}

var ErrUnauthorized = errors.New("unauthorized action")
var ErrInvalidpreftype = errors.New("invalid preference type")

type AIUseCase interface {
	GenerateIntialSuggestion(ctx context.Context, title string) (string, error)
	GenerateBasedOnTags(ctx context.Context, content string, tags []string) (string, error)
}
