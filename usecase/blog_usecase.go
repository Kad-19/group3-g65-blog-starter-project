package usecase

import (
    "context"
    "g3-g65-bsp/domain"
)


type blogUsecase struct {
    repo domain.BlogRepository
}

func NewBlogUsecase(repo domain.BlogRepository) domain.BlogUsecase {
    return &blogUsecase{repo: repo}
}

func (u *blogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog) (string, error) {
    return u.repo.CreateBlog(ctx, blog)
}

func (u *blogUsecase) GetBlogByID(ctx context.Context, id string) (*domain.Blog, error) {
    return u.repo.GetBlogByID(ctx, id)
}

func (u *blogUsecase) UpdateBlog(ctx context.Context, blog *domain.Blog) error {
    return u.repo.UpdateBlog(ctx, blog)
}

func (u *blogUsecase) DeleteBlog(ctx context.Context, id string) error {
    return u.repo.DeleteBlog(ctx, id)
}

func (u *blogUsecase) ListBlogs(ctx context.Context, filter map[string]interface{}) ([]*domain.Blog, error) {
    return u.repo.ListBlogs(ctx, filter)
}
