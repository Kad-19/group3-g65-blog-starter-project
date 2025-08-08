package usecase

import (
	"context"
	"g3-g65-bsp/domain"
	"strings"
)

type AIUsecaseImpl struct {
	aiclient domain.AIService
}

func NewAIUsecaseImpl(aus domain.AIService) *AIUsecaseImpl {
	return &AIUsecaseImpl{
		aiclient: aus,
	}
}

func (aus *AIUsecaseImpl) GenerateIntialSuggestion(ctx context.Context, title string) (string, error) {
	var promptBuilder strings.Builder
	promptBuilder.WriteString("Generate an engaging and informative piece of content based on the following title: ")
	promptBuilder.WriteString(title)
	promptBuilder.WriteString(". The content should be original, clear, and relevant to the topic.\n\n")
	promptBuilder.WriteString("The output must strictly contain the blog and nothing else.\n\n")

	res, err := aus.aiclient.GenerateContent(ctx, promptBuilder.String())
	if err != nil {
		return "", err
	}
	return res, nil
}

func (aus *AIUsecaseImpl) GenerateBasedOnTags(ctx context.Context, content string, tags []string) (string, error) {
	var promptBuilder strings.Builder
	promptBuilder.WriteString("Enhance the following content to be more engaging and informative. The enhancement should focus on these themes and keywords: ")
	promptBuilder.WriteString(strings.Join(tags, ", "))
	promptBuilder.WriteString("\n\nOriginal Content:\n")
	promptBuilder.WriteString(content)
	promptBuilder.WriteString("The output must strictly contain the blog and nothing else.\n\n")

	res, err := aus.aiclient.GenerateContent(ctx, promptBuilder.String())
	if err != nil {
		return "", err
	}
	return res, nil
}
