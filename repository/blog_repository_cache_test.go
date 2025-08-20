package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"g3-g65-bsp/domain"
)

// MockBlogRepository is a mock implementation of the BlogRepository interface for testing the cache.
type MockBlogRepository struct {
	mock.Mock
}

func (m *MockBlogRepository) CreateBlog(ctx context.Context, blog *domain.Blog) (string, error) {
	args := m.Called(ctx, blog)
	return args.String(0), args.Error(1)
}

func (m *MockBlogRepository) GetBlogByID(ctx context.Context, id string) (*domain.Blog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Blog), args.Error(1)
}

func (m *MockBlogRepository) UpdateBlog(ctx context.Context, blog *domain.Blog) error {
	args := m.Called(ctx, blog)
	return args.Error(0)
}

func (m *MockBlogRepository) DeleteBlog(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBlogRepository) ListBlogs(ctx context.Context, filter map[string]any, page, limit int) ([]*domain.Blog, *domain.Pagination, error) {
	args := m.Called(ctx, filter, page, limit)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*domain.Blog), args.Get(1).(*domain.Pagination), args.Error(2)
}

func (m *MockBlogRepository) IncrementBlogViewCount(ctx context.Context, id string, blog *domain.Blog) error {
	args := m.Called(ctx, id, blog)
	return args.Error(0)
}

func (m *MockBlogRepository) AddComment(ctx context.Context, blogID string, comment *domain.Comment) error {
	args := m.Called(ctx, blogID, comment)
	return args.Error(0)
}

func (m *MockBlogRepository) UpdateComment(ctx context.Context, blogID string, comment *domain.Comment) error {
	args := m.Called(ctx, blogID, comment)
	return args.Error(0)
}

func (m *MockBlogRepository) GetCommentByID(ctx context.Context, blogID string, commentID string) (*domain.Comment, error) {
	args := m.Called(ctx, blogID, commentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockBlogRepository) DeleteComment(ctx context.Context, blogID string, commentID string) error {
	args := m.Called(ctx, blogID, commentID)
	return args.Error(0)
}

// MockCacheService is a mock implementation of the CacheService interface.
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Set(key string, value interface{}, duration time.Duration) {
	m.Called(key, value, duration)
}

func (m *MockCacheService) Get(key string) (interface{}, bool) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0), args.Bool(1)
}

func (m *MockCacheService) Delete(key string) {
	m.Called(key)
}

func TestCachedBlogRepository_GetBlogByID(t *testing.T) {
	mockRepo := new(MockBlogRepository)
	mockCache := new(MockCacheService)
	cachedRepo := &cachedBlogRepository{
		repo:       mockRepo,
		cache:      mockCache,
		defaultTTL: 10 * time.Minute,
	}

	ctx := context.Background()
	blogID := "blog123"
	blog := &domain.Blog{ID: blogID, Title: "Test Blog"}

	// Test case 1: Cache hit
	mockCache.On("Get", blogCacheKey(blogID)).Return(blog, true).Once()
	result, err := cachedRepo.GetBlogByID(ctx, blogID)
	assert.NoError(t, err)
	assert.Equal(t, blog, result)
	mockCache.AssertExpectations(t)

	// Test case 2: Cache miss, get from repo and set cache
	mockCache.On("Get", blogCacheKey(blogID)).Return(nil, false).Once()
	mockRepo.On("GetBlogByID", ctx, blogID).Return(blog, nil).Once()
	mockCache.On("Set", blogCacheKey(blogID), blog, 10*time.Minute).Return().Once()
	result, err = cachedRepo.GetBlogByID(ctx, blogID)
	assert.NoError(t, err)
	assert.Equal(t, blog, result)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}