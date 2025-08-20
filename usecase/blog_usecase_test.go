package usecase

import (
	"context"
	"g3-g65-bsp/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBlogRepository is a mock implementation of the BlogRepository interface.
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

func (m *MockBlogRepository) ListBlogs(ctx context.Context, filter map[string]any, page, limit int) ([]*domain.Blog, *domain.Pagination, error) {
	args := m.Called(ctx, filter, page, limit)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*domain.Blog), args.Get(1).(*domain.Pagination), args.Error(2)
}

func TestBlogUsecase_CreateBlog(t *testing.T) {
	mockBlogRepo := new(MockBlogRepository)
	mockUserRepo := new(MockUserRepository)
	uc := NewBlogUsecase(mockBlogRepo, mockUserRepo)

	ctx := context.Background()
	userID := "user123"
	blog := &domain.Blog{
		Title:   "Test Title",
		Content: "Test Content",
	}

	// Mock user repository
	mockUserRepo.On("FindByID", ctx, userID).Return(&domain.User{Username: "testuser"}, nil).Once()
	// Mock blog repository
	mockBlogRepo.On("CreateBlog", ctx, mock.AnythingOfType("*domain.Blog")).Return("blog123", nil).Once()

	createdBlog, err := uc.CreateBlog(ctx, blog, userID)

	assert.NoError(t, err)
	assert.NotNil(t, createdBlog)
	assert.Equal(t, "blog123", createdBlog.ID)
	assert.Equal(t, userID, createdBlog.AuthorID)
	assert.Equal(t, "testuser", createdBlog.AuthorUsername)
	mockBlogRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestBlogUsecase_GetBlogByID(t *testing.T) {
	mockBlogRepo := new(MockBlogRepository)
	uc := NewBlogUsecase(mockBlogRepo, nil)

	ctx := context.Background()
	blogID := "blog123"
	expectedBlog := &domain.Blog{
		ID:      blogID,
		Title:   "Test Title",
		Content: "Test Content",
		Metrics: &domain.Metrics{ViewCount: 0},
	}

	// Mock blog repository
	mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(expectedBlog, nil).Once()
	mockBlogRepo.On("IncrementBlogViewCount", mock.Anything, blogID, mock.AnythingOfType("*domain.Blog")).Return(nil).Once()

	blog, err := uc.GetBlogByID(ctx, blogID)
	time.Sleep(50 * time.Millisecond) // allow goroutine to execute

	assert.NoError(t, err)
	assert.NotNil(t, blog)
	assert.Equal(t, expectedBlog, blog)
	assert.Equal(t, 1, blog.Metrics.ViewCount) // Check if view count was incremented
	mockBlogRepo.AssertExpectations(t)
}

func TestBlogUsecase_UpdateBlog(t *testing.T) {
	mockBlogRepo := new(MockBlogRepository)
	uc := NewBlogUsecase(mockBlogRepo, nil)

	ctx := context.Background()
	userID := "user123"
	blogID := "blog123"
	existingBlog := &domain.Blog{
		ID:       blogID,
		AuthorID: userID,
		Title:    "Old Title",
		Content:  "Old Content",
	}
	updatedBlog := &domain.Blog{
		Title:   "New Title",
		Content: "New Content",
	}

	// Mock blog repository
	mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(existingBlog, nil).Once()
	mockBlogRepo.On("UpdateBlog", ctx, mock.AnythingOfType("*domain.Blog")).Return(nil).Once()

	result, err := uc.UpdateBlog(ctx, updatedBlog, userID, blogID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "New Content", result.Content)
	mockBlogRepo.AssertExpectations(t)
}

func TestBlogUsecase_DeleteBlog(t *testing.T) {
	mockBlogRepo := new(MockBlogRepository)
	uc := NewBlogUsecase(mockBlogRepo, nil)

	ctx := context.Background()
	userID := "user123"
	adminID := "admin456"
	blogID := "blog123"
	existingBlog := &domain.Blog{
		ID:       blogID,
		AuthorID: userID,
	}

	// Test case 1: Successful deletion by author
	mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(existingBlog, nil).Once()
	mockBlogRepo.On("DeleteBlog", ctx, blogID).Return(nil).Once()
	err := uc.DeleteBlog(ctx, blogID, userID, "user")
	assert.NoError(t, err)

	// Test case 2: Successful deletion by admin
	mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(existingBlog, nil).Once()
	mockBlogRepo.On("DeleteBlog", ctx, blogID).Return(nil).Once()
	err = uc.DeleteBlog(ctx, blogID, adminID, "admin")
	assert.NoError(t, err)

	// Test case 3: Unauthorized deletion
	mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(existingBlog, nil).Once()
	err = uc.DeleteBlog(ctx, blogID, "anotheruser", "user")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUnauthorized, err)

	mockBlogRepo.AssertExpectations(t)
}
