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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockInteractionUsecase struct {
	mock.Mock
}

func (m *MockInteractionUsecase) LikeBlog(ctx context.Context, userID, blogID, prefType string) error {
	args := m.Called(ctx, userID, blogID, prefType)
	return args.Error(0)
}

func (m *MockInteractionUsecase) CommentOnBlog(ctx context.Context, userID, blogID string, comment *domain.Comment) error {
	args := m.Called(ctx, userID, blogID, comment)
	return args.Error(0)
}

func (m *MockInteractionUsecase) UpdateComment(ctx context.Context, userID, blogID, commentID, content string) error {
	args := m.Called(ctx, userID, blogID, commentID, content)
	return args.Error(0)
}

func (m *MockInteractionUsecase) DeleteComment(ctx context.Context, userID, blogID, commentID string) error {
	args := m.Called(ctx, userID, blogID, commentID)
	return args.Error(0)
}

func TestInteractionController_LikeBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockInteractionUsecase := new(MockInteractionUsecase)
		interactionController := NewInteractionController(mockInteractionUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "blog123"}}

		reqBody := LikeRequest{Preftype: "like"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/blogs/blog123/like", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockInteractionUsecase.On("LikeBlog", mock.Anything, "user123", "blog123", "like").Return(nil)

		interactionController.LikeBlog(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockInteractionUsecase.AssertExpectations(t)
	})
}

func TestInteractionController_CommentOnBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockInteractionUsecase := new(MockInteractionUsecase)
		interactionController := NewInteractionController(mockInteractionUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "blog123"}}

		reqBody := CommentRequest{Content: "Test Comment"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/blogs/blog123/comments", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockInteractionUsecase.On("CommentOnBlog", mock.Anything, "user123", "blog123", mock.AnythingOfType("*domain.Comment")).Return(nil)

		interactionController.CommentOnBlog(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockInteractionUsecase.AssertExpectations(t)
	})
}

func TestInteractionController_UpdateComment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockInteractionUsecase := new(MockInteractionUsecase)
		interactionController := NewInteractionController(mockInteractionUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "blog123"}, gin.Param{Key: "comment_id", Value: "comment123"}}

		reqBody := CommentRequest{Content: "Updated Comment"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/blogs/blog123/comments/comment123", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockInteractionUsecase.On("UpdateComment", mock.Anything, "user123", "blog123", "comment123", "Updated Comment").Return(nil)

		interactionController.UpdateComment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockInteractionUsecase.AssertExpectations(t)
	})
}

func TestInteractionController_DeleteComment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockInteractionUsecase := new(MockInteractionUsecase)
		interactionController := NewInteractionController(mockInteractionUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "blog123"}, gin.Param{Key: "comment_id", Value: "comment123"}}
		c.Request, _ = http.NewRequest(http.MethodDelete, "/blogs/blog123/comments/comment123", nil)

		mockInteractionUsecase.On("DeleteComment", mock.Anything, "user123", "blog123", "comment123").Return(nil)

		interactionController.DeleteComment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockInteractionUsecase.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockInteractionUsecase := new(MockInteractionUsecase)
		interactionController := NewInteractionController(mockInteractionUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "blog123"}, gin.Param{Key: "comment_id", Value: "comment123"}}
		c.Request, _ = http.NewRequest(http.MethodDelete, "/blogs/blog123/comments/comment123", nil)

		mockInteractionUsecase.On("DeleteComment", mock.Anything, "user123", "blog123", "comment123").Return(errors.New("some error"))

		interactionController.DeleteComment(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockInteractionUsecase.AssertExpectations(t)
	})
}
