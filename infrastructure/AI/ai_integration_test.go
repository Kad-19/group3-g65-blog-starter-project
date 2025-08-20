package ai

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGenerativeModel struct {
	mock.Mock
}

func (m *MockGenerativeModel) GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
	args := m.Called(ctx, parts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*genai.GenerateContentResponse), args.Error(1)
}

func TestNewGeminiService(t *testing.T) {
	// Test case 1: API key is set
	os.Setenv("GEMINI_AI_API_KEY", "test-api-key")
	service := NewGeminiService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.client)

	// Test case 2: API key is not set (should panic)
	os.Unsetenv("GEMINI_AI_API_KEY")
	assert.Panics(t, func() {
		NewGeminiService()
	})
}

func TestGenerateContent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockClient := new(MockGenerativeModel)
		service := &geminiService{client: mockClient}

		prompt := "test prompt"
		expectedResponse := "test response"
		apiResponse := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{
				{
					Content: &genai.Content{
						Parts: []genai.Part{genai.Text(expectedResponse)},
					},
				},
			},
		}

		mockClient.On("GenerateContent", mock.Anything, mock.Anything).Return(apiResponse, nil)

		content, err := service.GenerateContent(context.Background(), prompt)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, content)
		mockClient.AssertExpectations(t)
	})

	t.Run("api error", func(t *testing.T) {
		mockClient := new(MockGenerativeModel)
		service := &geminiService{client: mockClient}

		prompt := "test prompt"
		expectedError := errors.New("api error")

		mockClient.On("GenerateContent", mock.Anything, mock.Anything).Return(nil, expectedError)

		content, err := service.GenerateContent(context.Background(), prompt)

		assert.Error(t, err)
		assert.Empty(t, content)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockClient.AssertExpectations(t)
	})

	t.Run("malformed response", func(t *testing.T) {
		mockClient := new(MockGenerativeModel)
		service := &geminiService{client: mockClient}

		prompt := "test prompt"
		apiResponse := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{
				{
					Content: &genai.Content{
						Parts: []genai.Part{}, // Empty parts
					},
				},
			},
		}

		mockClient.On("GenerateContent", mock.Anything, mock.Anything).Return(apiResponse, nil)

		content, err := service.GenerateContent(context.Background(), prompt)

		assert.Error(t, err)
		assert.Empty(t, content)
		assert.Equal(t, "gemini api response was empty or malformed", err.Error())
		mockClient.AssertExpectations(t)
	})
}
