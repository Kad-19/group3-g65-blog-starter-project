package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogRepository interface {
	CreateBlog(ctx context.Context, blog *Blog) (string, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	UpdateBlog(ctx context.Context, blog *Blog) error
	DeleteBlog(ctx context.Context, id string) error
	ListBlogs(ctx context.Context, filter map[string]interface{}) ([]*Blog, error)
}

type UserRepositorys interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	UpdateUser(ctx context.Context, up UserProfile, Email string) error
	UpdateUserRole(ctx context.Context, role string, Email string) error
}
