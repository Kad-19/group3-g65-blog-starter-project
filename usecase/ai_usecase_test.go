package usecase

import (
	"context"
	"errors"
	"g3-g65-bsp/infrastructure"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	infrastructure.Log = log.New(io.Discard, "TEST: ", log.LstdFlags)
	os.Exit(m.Run())
}

// MockAIService is a mock of the domain.AIService interface
type MockAIService struct {
	mock.Mock
}

func (m *MockAIService) GenerateContent(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

func TestAIUsecaseImpl_GenerateIntialSuggestion(t *testing.T) {
	mockAIService := new(MockAIService)
	aiUsecase := NewAIUsecaseImpl(mockAIService)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedContent := "This is the generated content."
		jsonResponse := `{"appropriate": true, "description": "", "content": "This is the generated content."}`
		mockAIService.On("GenerateContent", ctx, mock.AnythingOfType("string")).Return(jsonResponse, nil).Once()

		content, err := aiUsecase.GenerateIntialSuggestion(ctx, "A good title")

		assert.NoError(t, err)
		assert.Equal(t, expectedContent, content)
		mockAIService.AssertExpectations(t)
	})

	t.Run("inappropriate", func(t *testing.T) {
		expectedDescription := "Inappropriate title"
		jsonResponse := `{"appropriate": false, "description": "Inappropriate title", "content": ""}`
		mockAIService.On("GenerateContent", ctx, mock.AnythingOfType("string")).Return(jsonResponse, nil).Once()

		content, err := aiUsecase.GenerateIntialSuggestion(ctx, "A bad title")

		assert.NoError(t, err)
		assert.Equal(t, expectedDescription, content)
		mockAIService.AssertExpectations(t)
	})

	t.Run("ai error", func(t *testing.T) {
		mockAIService.On("GenerateContent", ctx, mock.AnythingOfType("string")).Return("", errors.New("AI error")).Once()

		_, err := aiUsecase.GenerateIntialSuggestion(ctx, "Any title")

		assert.Error(t, err)
		mockAIService.AssertExpectations(t)
	})

	t.Run("json error", func(t *testing.T) {
		mockAIService.On("GenerateContent", ctx, mock.AnythingOfType("string")).Return("not a json", nil).Once()

		_, err := aiUsecase.GenerateIntialSuggestion(ctx, "Any title")

		assert.Error(t, err)
		mockAIService.AssertExpectations(t)
	})
}

func TestAIUsecaseImpl_GenerateBasedOnTags(t *testing.T) {
	mockAIService := new(MockAIService)
	aiUsecase := NewAIUsecaseImpl(mockAIService)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedContent := "Enhanced content."
		jsonResponse := `{"appropriate": true, "description": "", "content": "Enhanced content."}`
		mockAIService.On("GenerateContent", ctx, mock.AnythingOfType("string")).Return(jsonResponse, nil).Once()

		content, err := aiUsecase.GenerateBasedOnTags(ctx, "Original", []string{"tech"})

		assert.NoError(t, err)
		assert.Equal(t, expectedContent, content)
		mockAIService.AssertExpectations(t)
	})

	t.Run("inappropriate", func(t *testing.T) {
		expectedDescription := "Inappropriate content"
		jsonResponse := `{"appropriate": false, "description": "Inappropriate content", "content": ""}`
		mockAIService.On("GenerateContent", ctx, mock.AnythingOfType("string")).Return(jsonResponse, nil).Once()

		content, err := aiUsecase.GenerateBasedOnTags(ctx, "Bad content", []string{"tech"})

		assert.NoError(t, err)
		assert.Equal(t, expectedDescription, content)
		mockAIService.AssertExpectations(t)
	})

	t.Run("ai error", func(t *testing.T) {
		mockAIService.On("GenerateContent", ctx, mock.AnythingOfType("string")).Return("", errors.New("AI error")).Once()

		_, err := aiUsecase.GenerateBasedOnTags(ctx, "Any content", []string{"tech"})

		assert.Error(t, err)
		mockAIService.AssertExpectations(t)
	})

	t.Run("json error", func(t *testing.T) {
		mockAIService.On("GenerateContent", ctx, mock.AnythingOfType("string")).Return("not a json", nil).Once()

		_, err := aiUsecase.GenerateBasedOnTags(ctx, "Any content", []string{"tech"})

		assert.Error(t, err)
		mockAIService.AssertExpectations(t)
	})
}