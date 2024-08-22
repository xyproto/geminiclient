package simplegemini

import (
	"context"

	"cloud.google.com/go/vertexai/genai"
)

// CountPromptTokensWithClient counts the tokens in the given text prompt using a specific client and model.
func (gc *GeminiClient) CountPromptTokensWithClient(ctx context.Context, client *genai.Client, prompt, modelName string) (int, error) {
	model := client.GenerativeModel(modelName)
	resp, err := model.CountTokens(ctx, genai.Text(prompt))
	if err != nil {
		return 0, err
	}
	return int(resp.TotalTokens), nil
}

// CountPromptTokensWithModel counts the tokens in the given text prompt using the specified model within the default client.
func (gc *GeminiClient) CountPromptTokensWithModel(ctx context.Context, prompt, modelName string) (int, error) {
	model := gc.Client.GenerativeModel(modelName)
	resp, err := model.CountTokens(ctx, genai.Text(prompt))
	if err != nil {
		return 0, err
	}
	return int(resp.TotalTokens), nil
}

// CountPromptTokens counts the number of tokens in the given text prompt using the default client and model.
func (gc *GeminiClient) CountPromptTokens(prompt string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gc.Timeout)
	defer cancel()
	return gc.CountPromptTokensWithModel(ctx, prompt, gc.ModelName)
}

// CountPartTokensWithContext counts the tokens in the current multimodal parts using the default client and model.
func (gc *GeminiClient) CountPartTokensWithContext(ctx context.Context) (int, error) {
	model := gc.Client.GenerativeModel(gc.ModelName)
	var totalTokens int
	for _, part := range gc.Parts {
		resp, err := model.CountTokens(ctx, part)
		if err != nil {
			return totalTokens, err
		}
		totalTokens += int(resp.TotalTokens)
	}
	return totalTokens, nil
}

// CountPartTokens counts the tokens in the current multimodal parts using the default client, model, and a new context.
func (gc *GeminiClient) CountPartTokens() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gc.Timeout)
	defer cancel()
	return gc.CountPartTokensWithContext(ctx)
}

// CountTextTokensWithClient counts the tokens in the given text using a specific client and model.
func (gc *GeminiClient) CountTextTokensWithClient(ctx context.Context, client *genai.Client, text, modelName string) (int, error) {
	model := client.GenerativeModel(modelName)
	resp, err := model.CountTokens(ctx, genai.Text(text))
	if err != nil {
		return 0, err
	}
	return int(resp.TotalTokens), nil
}

// CountTextTokensWithModel counts the tokens in the given text using the specified model within the default client.
func (gc *GeminiClient) CountTextTokensWithModel(ctx context.Context, text, modelName string) (int, error) {
	model := gc.Client.GenerativeModel(modelName)
	resp, err := model.CountTokens(ctx, genai.Text(text))
	if err != nil {
		return 0, err
	}
	return int(resp.TotalTokens), nil
}

// CountTextTokens counts the tokens in the given text using the default client and model.
func (gc *GeminiClient) CountTextTokens(text string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gc.Timeout)
	defer cancel()
	return gc.CountTextTokensWithModel(ctx, text, gc.ModelName)
}
