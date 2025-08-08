package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure"
	"strings"
)

type AIUsecaseImpl struct {
	aiclient domain.AIService
}

type EvaluationResult struct {
	Appropriate bool
	Description string
	Content     string
}

func NewAIUsecaseImpl(aus domain.AIService) *AIUsecaseImpl {
	return &AIUsecaseImpl{
		aiclient: aus,
	}
}

func (aus *AIUsecaseImpl) GenerateIntialSuggestion(ctx context.Context, title string) (string, error) {
	var promptBuilder strings.Builder
	promptBuilder.WriteString(`You are an AI content evaluator and generator.`)
	promptBuilder.WriteString(` Evaluate the following blog title and determine if it's appropriate for generating blog content.`)
	promptBuilder.WriteString(` If it's appropriate, generate engaging and informative content based on it.`)
	promptBuilder.WriteString(` If not appropriate, provide a reason.\n\n`)
	promptBuilder.WriteString(fmt.Sprintf("Title: \"%s\"\n\n", title))

	promptBuilder.WriteString(`Respond strictly in the following JSON format:
{
  "appropriate": true/false,
  "description": "Reason if inappropriate, otherwise empty",
  "content": "Generated content if appropriate, otherwise empty"
}
`)

	res, err := aus.aiclient.GenerateContent(ctx, promptBuilder.String())
	if err != nil {
		return "", err
	}

	cleanRes := strings.TrimSpace(res)
	cleanRes = strings.TrimPrefix(cleanRes, "```json")
	cleanRes = strings.TrimPrefix(cleanRes, "```")
	cleanRes = strings.TrimSuffix(cleanRes, "```")
	cleanRes = strings.TrimSpace(cleanRes)

	var eval EvaluationResult
	if err := json.Unmarshal([]byte(cleanRes), &eval); err != nil {
		return "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	// infrastructure.Log.Println(eval)
	if eval.Appropriate {
		return eval.Content, nil
	}

	return eval.Description, nil
}

func (aus *AIUsecaseImpl) GenerateBasedOnTags(ctx context.Context, content string, tags []string) (string, error) {
	var promptBuilder strings.Builder
	promptBuilder.WriteString(`You are an AI content evaluator and enhancer.`)
	promptBuilder.WriteString(` Evaluate the content and determine whether it's appropriate to enhance using the given tags.`)
	promptBuilder.WriteString(` If appropriate, enhance the content based on the tags to improve its quality and relevance.`)
	promptBuilder.WriteString(` If not, explain why enhancement is not suitable.\n\n`)

	promptBuilder.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(tags, ", ")))
	promptBuilder.WriteString("Original Content:\n")
	promptBuilder.WriteString(content)

	promptBuilder.WriteString(`
Respond with ONLY valid JSON. 
Do not include markdown code blocks, backticks, or any other formatting — only the JSON object.
`)

	// 	promptBuilder.WriteString(`\n\nRespond strictly in the following JSON format:
	// {
	//   "appropriate": true/false,
	//   "description": "Reason if inappropriate, otherwise empty",
	//   "content": "Enhanced content if appropriate, otherwise empty"
	// }
	// `)
	res, err := aus.aiclient.GenerateContent(ctx, promptBuilder.String())
	if err != nil {
		return "", err
	}

	cleanRes := strings.TrimSpace(res)
	cleanRes = strings.TrimPrefix(cleanRes, "```json")
	cleanRes = strings.TrimPrefix(cleanRes, "```")
	cleanRes = strings.TrimSuffix(cleanRes, "```")
	cleanRes = strings.TrimSpace(cleanRes)

	var eval EvaluationResult
	if err := json.Unmarshal([]byte(cleanRes), &eval); err != nil {
		return "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	infrastructure.Log.Println(eval)
	if eval.Appropriate {
		return eval.Content, nil
	}

	return eval.Description, nil
}
