package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAIUseCase struct {
	mock.Mock
}

func (m *MockAIUseCase) GenerateIntialSuggestion(ctx context.Context, title string) (string, error) {
	args := m.Called(ctx, title)
	return args.String(0), args.Error(1)
}

func (m *MockAIUseCase) GenerateBasedOnTags(ctx context.Context, content string, tags []string) (string, error) {
	args := m.Called(ctx, content, tags)
	return args.String(0), args.Error(1)
}

func (m *MockAIUseCase) GenerateSummary(ctx context.Context, content string) (string, error) {
	args := m.Called(ctx, content)
	return args.String(0), args.Error(1)
}

func TestAIcontroller_HandleAIContentrequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockAIUseCase := new(MockAIUseCase)
		aiController := NewAIcontroller(mockAIUseCase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := ContentRequest{Title: "Test Title"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		expectedResponse := "Generated Content"
		mockAIUseCase.On("GenerateIntialSuggestion", mock.Anything, reqBody.Title).Return(expectedResponse, nil)

		aiController.HandleAIContentrequest(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseBody map[string]string
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, expectedResponse, responseBody["content"])
		mockAIUseCase.AssertExpectations(t)
	})

	t.Run("bad request", func(t *testing.T) {
		aiController := NewAIcontroller(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("{")))
		c.Request.Header.Set("Content-Type", "application/json")

		aiController.HandleAIContentrequest(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockAIUseCase := new(MockAIUseCase)
		aiController := NewAIcontroller(mockAIUseCase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := ContentRequest{Title: "Test Title"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockAIUseCase.On("GenerateIntialSuggestion", mock.Anything, reqBody.Title).Return("", errors.New("some error"))

		aiController.HandleAIContentrequest(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockAIUseCase.AssertExpectations(t)
	})
}

func TestAIcontroller_HandleAIEnhancement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockAIUseCase := new(MockAIUseCase)
		aiController := NewAIcontroller(mockAIUseCase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := EnhanceContent{Content: "Test Content", Tags: []string{"tag1", "tag2"}}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		expectedResponse := "Enhanced Content"
		mockAIUseCase.On("GenerateBasedOnTags", mock.Anything, reqBody.Content, reqBody.Tags).Return(expectedResponse, nil)

		aiController.HandleAIEnhancement(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseBody map[string]string
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, expectedResponse, responseBody["content"])
		mockAIUseCase.AssertExpectations(t)
	})

	t.Run("bad request", func(t *testing.T) {
		mockAIUseCase := new(MockAIUseCase)
		aiController := NewAIcontroller(mockAIUseCase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("")))
		c.Request.Header.Set("Content-Type", "application/json")

		aiController.HandleAIEnhancement(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockAIUseCase := new(MockAIUseCase)
		aiController := NewAIcontroller(mockAIUseCase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := EnhanceContent{Content: "Test Content", Tags: []string{"tag1", "tag2"}}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockAIUseCase.On("GenerateBasedOnTags", mock.Anything, reqBody.Content, reqBody.Tags).Return("", errors.New("some error"))

		aiController.HandleAIEnhancement(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockAIUseCase.AssertExpectations(t)
	})
}
