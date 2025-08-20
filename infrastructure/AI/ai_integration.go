package ai

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GenerativeModelAPI interface {
	GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
}

type geminiService struct {
	client GenerativeModelAPI
}

func NewGeminiService() *geminiService {
	apikey := os.Getenv("GEMINI_AI_API_KEY")
	if apikey == "" {
		panic("GOOGLE_API_KEY environment variable not set")
	}
	ctx := context.Background()
	cld, err := genai.NewClient(ctx, option.WithAPIKey(apikey))
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Gemini client: %v", err))
	}
	model := cld.GenerativeModel("gemini-2.5-flash")
	return &geminiService{
		client: model,
	}
}

func (gs *geminiService) GenerateContent(ctx context.Context, prompt string) (string, error) {
	resp, err := gs.client.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content from gemini api: %w", err)
	}
	if resp != nil && len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
		if ok {
			return string(textPart), nil
		}
	}

	return "", errors.New("gemini api response was empty or malformed")
}
