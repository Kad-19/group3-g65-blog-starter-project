package domain

import (
	"context"
)

type BlogRepository interface {
	CreateBlog(ctx context.Context, blog *Blog) (string, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	UpdateBlog(ctx context.Context, blog *Blog) error
	DeleteBlog(ctx context.Context, id string) error
	ListBlogs(ctx context.Context, filter map[string]any, page, limit int) ([]*Blog, *Pagination, error)
	IncrementBlogViewCount(ctx context.Context, id string, blog *Blog) error
	AddComment(ctx context.Context, blogID string, comment *Comment) error
	UpdateComment(ctx context.Context, blogID string, comment *Comment) error
	GetCommentByID(ctx context.Context, blogID string, commentID string) (*Comment, error)
	DeleteComment(ctx context.Context, blogID string, commentID string) error
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	UpdateUserProfile(ctx context.Context, bio string, contactInfo string, imagePath string, Email string) error
	UpdateUserRole(ctx context.Context, role string, Email string) error
	UpdateActiveStatus(ctx context.Context, email string) error
	UpdateUserPassword(ctx context.Context, email string, newPasswordHash string) error
	GetAllUsers(ctx context.Context, page int, limit int) ([]User, int64, error)
}

type UnactiveUserRepo interface {
	CreateUnactiveUser(ctx context.Context, user *UnactivatedUser) error
	FindByEmailUnactive(ctx context.Context, email string) (*UnactivatedUser, error)
	DeleteUnactiveUser(ctx context.Context, email string) error
}

type PasswordResetRepository interface {
	Create(ctx context.Context, token *PasswordResetToken) error
	GetByToken(ctx context.Context, token string) (*PasswordResetToken, error)
	Delete(ctx context.Context, token string) error
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, accessToken *RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) (error)
	DeleteAllForUser(ctx context.Context, userID string) error
}

type InteractionRepository interface {
	LikeBlog(ctx context.Context, userID string, blogID string, preftype string) error
	CommentOnBlog(ctx context.Context, userID string, blogID string, comment *Comment) error
}
