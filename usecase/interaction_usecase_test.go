package usecase

import (
	"context"
	"g3-g65-bsp/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInteractionRepository is a mock implementation of the InteractionRepository interface.
type MockInteractionRepository struct {
	mock.Mock
}

func (m *MockInteractionRepository) AddComment(ctx context.Context, blogID string, comment *domain.Comment) error {
	args := m.Called(ctx, blogID, comment)
	return args.Error(0)
}

func (m *MockInteractionRepository) GetCommentByID(ctx context.Context, blogID, commentID string) (*domain.Comment, error) {
	args := m.Called(ctx, blogID, commentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockInteractionRepository) UpdateComment(ctx context.Context, blogID string, comment *domain.Comment) error {
	args := m.Called(ctx, blogID, comment)
	return args.Error(0)
}

func (m *MockInteractionRepository) DeleteComment(ctx context.Context, blogID, commentID string) error {
	args := m.Called(ctx, blogID, commentID)
	return args.Error(0)
}

func TestInteractionUsecase_LikeBlog(t *testing.T) {
	mockBlogRepo := new(MockBlogRepository)
	uc := NewInteractionUsecase(mockBlogRepo, nil)

	ctx := context.Background()
	userID := "user123"
	blogID := "blog123"

	t.Run("like a blog for the first time", func(t *testing.T) {
		blog := &domain.Blog{
			ID: blogID,
			Metrics: &domain.Metrics{
				Likes:    &domain.Likes{Count: 0, Users: []string{}},
				Dislikes: &domain.Likes{Count: 0, Users: []string{}},
			},
		}
		mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(blog, nil).Once()
		mockBlogRepo.On("UpdateBlog", ctx, mock.AnythingOfType("*domain.Blog")).Return(nil).Once()

		err := uc.LikeBlog(ctx, userID, blogID, "like")
		assert.NoError(t, err)
		assert.Equal(t, 1, blog.Metrics.Likes.Count)
		assert.Contains(t, blog.Metrics.Likes.Users, userID)
		mockBlogRepo.AssertExpectations(t)
	})

	t.Run("unlike a blog", func(t *testing.T) {
		blog := &domain.Blog{
			ID: blogID,
			Metrics: &domain.Metrics{
				Likes:    &domain.Likes{Count: 1, Users: []string{userID}},
				Dislikes: &domain.Likes{Count: 0, Users: []string{}},
			},
		}
		mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(blog, nil).Once()
		mockBlogRepo.On("UpdateBlog", ctx, mock.AnythingOfType("*domain.Blog")).Return(nil).Once()

		err := uc.LikeBlog(ctx, userID, blogID, "like")
		assert.NoError(t, err)
		assert.Equal(t, 0, blog.Metrics.Likes.Count)
		assert.NotContains(t, blog.Metrics.Likes.Users, userID)
		mockBlogRepo.AssertExpectations(t)
	})

	t.Run("dislike a blog", func(t *testing.T) {
		blog := &domain.Blog{
			ID: blogID,
			Metrics: &domain.Metrics{
				Likes:    &domain.Likes{Count: 0, Users: []string{}},
				Dislikes: &domain.Likes{Count: 0, Users: []string{}},
			},
		}
		mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(blog, nil).Once()
		mockBlogRepo.On("UpdateBlog", ctx, mock.AnythingOfType("*domain.Blog")).Return(nil).Once()

		err := uc.LikeBlog(ctx, userID, blogID, "dislike")
		assert.NoError(t, err)
		assert.Equal(t, 1, blog.Metrics.Dislikes.Count)
		assert.Contains(t, blog.Metrics.Dislikes.Users, userID)
		mockBlogRepo.AssertExpectations(t)
	})

	t.Run("undislike a blog", func(t *testing.T) {
		blog := &domain.Blog{
			ID: blogID,
			Metrics: &domain.Metrics{
				Likes:    &domain.Likes{Count: 0, Users: []string{}},
				Dislikes: &domain.Likes{Count: 1, Users: []string{userID}},
			},
		}
		mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(blog, nil).Once()
		mockBlogRepo.On("UpdateBlog", ctx, mock.AnythingOfType("*domain.Blog")).Return(nil).Once()

		err := uc.LikeBlog(ctx, userID, blogID, "dislike")
		assert.NoError(t, err)
		assert.Equal(t, 0, blog.Metrics.Dislikes.Count)
		assert.NotContains(t, blog.Metrics.Dislikes.Users, userID)
		mockBlogRepo.AssertExpectations(t)
	})

	t.Run("change from like to dislike", func(t *testing.T) {
		blog := &domain.Blog{
			ID: blogID,
			Metrics: &domain.Metrics{
				Likes:    &domain.Likes{Count: 1, Users: []string{userID}},
				Dislikes: &domain.Likes{Count: 0, Users: []string{}},
			},
		}
		mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(blog, nil).Once()
		mockBlogRepo.On("UpdateBlog", ctx, mock.AnythingOfType("*domain.Blog")).Return(nil).Once()

		err := uc.LikeBlog(ctx, userID, blogID, "dislike")
		assert.NoError(t, err)
		assert.Equal(t, 0, blog.Metrics.Likes.Count)
		assert.Equal(t, 1, blog.Metrics.Dislikes.Count)
		assert.NotContains(t, blog.Metrics.Likes.Users, userID)
		assert.Contains(t, blog.Metrics.Dislikes.Users, userID)
		mockBlogRepo.AssertExpectations(t)
	})
}

func TestInteractionUsecase_CommentOnBlog(t *testing.T) {
	mockBlogRepo := new(MockBlogRepository)
	mockUserRepo := new(MockUserRepository)
	uc := NewInteractionUsecase(mockBlogRepo, mockUserRepo)

	ctx := context.Background()
	userID := "user123"
	blogID := "blog123"
	comment := &domain.Comment{Content: "Test comment"}

	mockUserRepo.On("FindByID", ctx, userID).Return(&domain.User{Username: "testuser"}, nil).Once()
	mockBlogRepo.On("GetBlogByID", ctx, blogID).Return(&domain.Blog{}, nil).Once()
	mockBlogRepo.On("AddComment", ctx, blogID, mock.AnythingOfType("*domain.Comment")).Return(nil).Once()

	err := uc.CommentOnBlog(ctx, userID, blogID, comment)

	assert.NoError(t, err)
	assert.Equal(t, userID, comment.AuthorID)
	assert.Equal(t, "testuser", comment.AuthorUsername)
	mockBlogRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// ... (add tests for UpdateComment and DeleteComment)
