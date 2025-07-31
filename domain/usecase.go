package domain

import "context"

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
