package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"g3-g65-bsp/domain"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBlogUsecase struct {
	mock.Mock
}

func (m *MockBlogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog, userID string) (*domain.Blog, error) {
	args := m.Called(ctx, blog, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Blog), args.Error(1)
}

func (m *MockBlogUsecase) GetBlogByID(ctx context.Context, id string) (*domain.Blog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Blog), args.Error(1)
}

func (m *MockBlogUsecase) UpdateBlog(ctx context.Context, blog *domain.Blog, userID, blogID string) (*domain.Blog, error) {
	args := m.Called(ctx, blog, userID, blogID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Blog), args.Error(1)
}

func (m *MockBlogUsecase) DeleteBlog(ctx context.Context, id, userID, role string) error {
	args := m.Called(ctx, id, userID, role)
	return args.Error(0)
}

func (m *MockBlogUsecase) ListBlogs(ctx context.Context, filter map[string]any, page, limit int) ([]*domain.Blog, *domain.Pagination, error) {
	args := m.Called(ctx, filter, page, limit)
	var blogs []*domain.Blog
	if args.Get(0) != nil {
		blogs = args.Get(0).([]*domain.Blog)
	}
	var pagination *domain.Pagination
	if args.Get(1) != nil {
		pagination = args.Get(1).(*domain.Pagination)
	}
	return blogs, pagination, args.Error(2)
}

func TestBlogController_CreateBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockBlogUsecase := new(MockBlogUsecase)
		blogController := NewBlogController(mockBlogUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		reqBody := BlogDTO{Title: "Test Title", Content: "Test Content"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/blogs", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		now := time.Now()
		blog := &domain.Blog{ID: "1", Title: "Test Title", Content: "Test Content", AuthorID: "user123", CreatedAt: &now, UpdatedAt: &now, Metrics: &domain.Metrics{Likes: &domain.Likes{}, Dislikes: &domain.Likes{}}}
		mockBlogUsecase.On("CreateBlog", mock.Anything, mock.AnythingOfType("*domain.Blog"), "user123").Return(blog, nil)

		blogController.CreateBlog(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response BlogDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Title", response.Title)
		mockBlogUsecase.AssertExpectations(t)
	})
}

func TestBlogController_GetBlogByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockBlogUsecase := new(MockBlogUsecase)
		blogController := NewBlogController(mockBlogUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		now := time.Now()
		blog := &domain.Blog{ID: "1", Title: "Test Title", Content: "Test Content", AuthorID: "user123", CreatedAt: &now, UpdatedAt: &now, Metrics: &domain.Metrics{Likes: &domain.Likes{}, Dislikes: &domain.Likes{}}}
		mockBlogUsecase.On("GetBlogByID", mock.Anything, "1").Return(blog, nil)

		blogController.GetBlogByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response BlogDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Title", response.Title)
		mockBlogUsecase.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockBlogUsecase := new(MockBlogUsecase)
		blogController := NewBlogController(mockBlogUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/blogs/1", nil)

		mockBlogUsecase.On("GetBlogByID", mock.Anything, "1").Return(nil, errors.New("not found"))

		blogController.GetBlogByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockBlogUsecase.AssertExpectations(t)
	})
}

func TestBlogController_UpdateBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockBlogUsecase := new(MockBlogUsecase)
		blogController := NewBlogController(mockBlogUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		reqBody := BlogDTO{Title: "Updated Title", Content: "Updated Content"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/blogs/1", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		now := time.Now()
		blog := &domain.Blog{ID: "1", Title: "Updated Title", Content: "Updated Content", AuthorID: "user123", CreatedAt: &now, UpdatedAt: &now, Metrics: &domain.Metrics{Likes: &domain.Likes{}, Dislikes: &domain.Likes{}}}
		mockBlogUsecase.On("UpdateBlog", mock.Anything, mock.AnythingOfType("*domain.Blog"), "user123", "1").Return(blog, nil)

		blogController.UpdateBlog(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response BlogDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Title", response.Title)
		mockBlogUsecase.AssertExpectations(t)
	})
}

func TestBlogController_DeleteBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockBlogUsecase := new(MockBlogUsecase)
		blogController := NewBlogController(mockBlogUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Set("role", "user")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest(http.MethodDelete, "/blogs/1", nil)

		mockBlogUsecase.On("DeleteBlog", mock.Anything, "1", "user123", "user").Return(nil)

		blogController.DeleteBlog(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockBlogUsecase.AssertExpectations(t)
	})
}

func TestBlogController_ListBlogs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockBlogUsecase)
	controller := NewBlogController(mockUsecase)

	t.Run("success", func(t *testing.T) {
		blogs := []*domain.Blog{{ID: "1", Title: "Test Blog", AuthorID: "user123", Content: "content", Tags: []string{"tag1"}, Metrics: &domain.Metrics{Likes: &domain.Likes{}, Dislikes: &domain.Likes{}}}}
		pagination := &domain.Pagination{Total: 1, Page: 1, Limit: 10}
		mockUsecase.On("ListBlogs", mock.Anything, mock.Anything, 1, 10).Return(blogs, pagination, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/blogs?page=1&limit=10", nil)

		controller.ListBlogs(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody struct {
			Data       []BlogDTO         `json:"data"`
			Pagination domain.Pagination `json:"pagination"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Len(t, responseBody.Data, 1)
		assert.Equal(t, "Test Blog", responseBody.Data[0].Title)
		assert.Equal(t, 1, responseBody.Pagination.Total)
	})
}
