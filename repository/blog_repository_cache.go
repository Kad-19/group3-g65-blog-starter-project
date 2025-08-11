package repository

import (
	"context"
	"errors"
	"fmt"
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CacheService defines the interface for a cache.
// This allows us to use any cache implementation (in-memory, Redis, etc.).
type CacheService interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

// cachedBlogRepository is a decorator that adds caching to a BlogRepository.
type cachedBlogRepository struct {
	repo       domain.BlogRepository // The "next" repository (our mongo implementation)
	cache      CacheService
	defaultTTL time.Duration
}

// ErrBlogNotFound is returned when a blog is not found in the repository
var ErrBlogNotFound = errors.New("blog not found")

// NewBlogRepository returns a MongoDB implementation of BlogRepository
func NewBlogRepository(collection *mongo.Collection, cache CacheService) domain.BlogRepository {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tags", Value: 1},
				{Key: "created_at", Value: -1},
				{Key: "metrics.view_count", Value: -1},
			},
		},
	}

	if _, err := collection.Indexes().CreateMany(context.Background(), indexes); err != nil {
		infrastructure.Log.Fatalf("Failed to create index: %v", err)
	}
	
	return &cachedBlogRepository{
		repo:       &mongoBlogRepository{collection: collection},
		cache:      cache,
		defaultTTL: 10 * time.Minute, // Set a default cache duration
	}
}

// blogCacheKey generates a unique cache key for a blog entry.
func blogCacheKey(id string) string {
	return fmt.Sprintf("blog:%s", id)
}

// GetBlogByID checks the cache first before hitting the database.
func (r *cachedBlogRepository) GetBlogByID(ctx context.Context, id string) (*domain.Blog, error) {
	key := blogCacheKey(id)

	// 1. Attempt to get the blog from the cache
	if cached, found := r.cache.Get(key); found {
		if blog, ok := cached.(*domain.Blog); ok {
			fmt.Print("Cache hit for blog ID: ", id, "\n")
			return blog, nil
		}
	}

	// 2. If not in cache, get from the underlying repository (database)
	blog, err := r.repo.GetBlogByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. Store the result in the cache for future requests
	r.cache.Set(key, blog, r.defaultTTL)

	return blog, nil
}

// NOTE: Caching ListBlogs is complex due to dynamic filters.
// For now, we will pass it through to the underlying repository.
func (r *cachedBlogRepository) ListBlogs(ctx context.Context, filter map[string]any, page, limit int) ([]*domain.Blog, *domain.Pagination, error) {
	return r.repo.ListBlogs(ctx, filter, page, limit)
}

// UpdateBlog updates the blog in the DB and invalidates the cache.
func (r *cachedBlogRepository) UpdateBlog(ctx context.Context, blog *domain.Blog) error {
	// First, execute the primary operation
	if err := r.repo.UpdateBlog(ctx, blog); err != nil {
		return err
	}

	// If successful, invalidate the cache
	r.cache.Delete(blogCacheKey(blog.ID))
	return nil
}

// DeleteBlog deletes the blog from the DB and invalidates the cache.
func (r *cachedBlogRepository) DeleteBlog(ctx context.Context, id string) error {
	if err := r.repo.DeleteBlog(ctx, id); err != nil {
		return err
	}
	r.cache.Delete(blogCacheKey(id))
	return nil
}

// AddComment invalidates the parent blog's cache.
func (r *cachedBlogRepository) AddComment(ctx context.Context, blogID string, comment *domain.Comment) error {
	if err := r.repo.AddComment(ctx, blogID, comment); err != nil {
		return err
	}
	r.cache.Delete(blogCacheKey(blogID))
	return nil
}

// UpdateComment invalidates the parent blog's cache.
func (r *cachedBlogRepository) UpdateComment(ctx context.Context, blogID string, comment *domain.Comment) error {
	if err := r.repo.UpdateComment(ctx, blogID, comment); err != nil {
		return err
	}
	r.cache.Delete(blogCacheKey(blogID))
	return nil
}

// DeleteComment invalidates the parent blog's cache.
func (r *cachedBlogRepository) DeleteComment(ctx context.Context, blogID string, commentID string) error {
	if err := r.repo.DeleteComment(ctx, blogID, commentID); err != nil {
		return err
	}
	r.cache.Delete(blogCacheKey(blogID))
	return nil
}

// IncrementBlogViewCount invalidates the blog's cache.
func (r *cachedBlogRepository) IncrementBlogViewCount(ctx context.Context, id string, blog *domain.Blog) error {

	return r.repo.IncrementBlogViewCount(ctx, id, blog)
}

// Pass-through methods that don't affect single-blog caching
func (r *cachedBlogRepository) CreateBlog(ctx context.Context, blog *domain.Blog) (string, error) {
	// No invalidation needed on create, but could optionally "warm" the cache
	return r.repo.CreateBlog(ctx, blog)
}

func (r *cachedBlogRepository) GetCommentByID(ctx context.Context, blogID string, commentID string) (*domain.Comment, error) {
	// Caching individual comments could be done, but for now we pass through
	return r.repo.GetCommentByID(ctx, blogID, commentID)
}
