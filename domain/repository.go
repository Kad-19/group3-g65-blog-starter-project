package domain

import "context"

type BlogRepository interface {
	CreateBlog(ctx context.Context, blog *Blog) (string, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	UpdateBlog(ctx context.Context, blog *Blog) error
	DeleteBlog(ctx context.Context, id string) error
	ListBlogs(ctx context.Context, filter map[string]interface{}) ([]*Blog, error)
}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, Email string) (*User, error)
	UpdateUser(ctx context.Context, up UserProfile, Email string) error
	UpdateUserRole(ctx context.Context, role string, Email string) error
}
