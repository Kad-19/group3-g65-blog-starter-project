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
    blog, err := u.repo.GetBlogByID(ctx, id)
    if err != nil {
        return nil, err
    }
    // Atomically increment the view count in the database
    if blog != nil {
        _ = u.repo.IncrementBlogViewCount(ctx, id, blog)

        blog.Metrics.ViewCount += 1 // reflect increment in returned object
    }
    return blog, nil
}

func (u *blogUsecase) UpdateBlog(ctx context.Context, blog *domain.Blog) error {
    return u.repo.UpdateBlog(ctx, blog)
}

func (u *blogUsecase) DeleteBlog(ctx context.Context, id string) error {
    return u.repo.DeleteBlog(ctx, id)
}

// ListBlogs allows filtering by tags ([]string), date (created_at_from, created_at_to), or popularity (min_views)
func (u *blogUsecase) ListBlogs(ctx context.Context, filter map[string]any, page, limit int) ([]*domain.Blog, *domain.Pagination, error) {
    return u.repo.ListBlogs(ctx, filter, page, limit)
}

