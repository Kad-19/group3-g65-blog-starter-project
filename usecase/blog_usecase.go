package usecase

import (
	"context"
	"g3-g65-bsp/domain"
	"time"

)


type blogUsecase struct {
    repo domain.BlogRepository
    userRepo domain.UserRepository
}

func NewBlogUsecase(repo domain.BlogRepository, userRepo domain.UserRepository) domain.BlogUsecase {
    return &blogUsecase{repo: repo, userRepo: userRepo}
}

func (u *blogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog, userid string) (*domain.Blog, error) {
    
    existingUser, err := u.userRepo.FindByID(ctx, userid)
    if err != nil {
        return nil, err
    }
    blog.AuthorID = userid
    blog.AuthorUsername = existingUser.Username


    now := time.Now()
    blog.CreatedAt = &now
    blog.UpdatedAt = blog.CreatedAt
    blog.Metrics = &domain.Metrics{
        ViewCount: 0,
        Likes:     &domain.Likes{Count: 0, Users: []string{}},
        Dislikes:  &domain.Likes{Count: 0, Users: []string{}},
    }
    blog.Comments = []domain.Comment{}


    blogid, err := u.repo.CreateBlog(ctx, blog)
    if err != nil {
        return nil, err
    }
    blog.ID = blogid
    return blog, nil
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

func (u *blogUsecase) UpdateBlog(ctx context.Context, blog *domain.Blog, userid, id string) (*domain.Blog, error) {
    // Ensure the blog belongs to the user
    existingBlog, err := u.repo.GetBlogByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if existingBlog.AuthorID != userid {
        return nil, domain.ErrUnauthorized
    }
    existingBlog.Title = blog.Title
    existingBlog.Content = blog.Content
    existingBlog.Tags = blog.Tags
    
    e := u.repo.UpdateBlog(ctx, existingBlog)
    if e != nil {
        return nil, e
    }
    return existingBlog, nil
}

func (u *blogUsecase) DeleteBlog(ctx context.Context, id string) error {
    return u.repo.DeleteBlog(ctx, id)
}

// ListBlogs allows filtering by tags ([]string), date (created_at_from, created_at_to), or popularity (min_views)
func (u *blogUsecase) ListBlogs(ctx context.Context, filter map[string]any, page, limit int) ([]*domain.Blog, *domain.Pagination, error) {
    return u.repo.ListBlogs(ctx, filter, page, limit)
}

