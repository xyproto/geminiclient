package simplegemini

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

// CountTextTokensWithClient will count the tokens in the given text.
func (gc *GeminiClient) CountTextTokensWithClient(ctx context.Context, client *genai.Client, text string) (int, error) {
	model := client.GenerativeModel(gc.ModelName)
	resp, err := model.CountTokens(ctx, genai.Text(text))
	if err != nil {
		return 0, err
	}
	return int(resp.TotalTokens), nil
}

// CountTextTokensWithModel will count the tokens in the given text and model
func (gc *GeminiClient) CountTextTokensWithModel(prompt, modelName string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gc.Timeout)
	defer cancel()

	model := gc.Client.GenerativeModel(modelName)
	resp, err := model.CountTokens(ctx, genai.Text(prompt))
	if err != nil {
		return 0, err
	}
	return int(resp.TotalTokens), nil
}

// CountTokensWithContext will count the tokens in the current multimodal prompt.
func (gc *GeminiClient) CountTokensWithContext(ctx context.Context) (int, error) {
	model := gc.Client.GenerativeModel(gc.ModelName)
	var sum int
	for _, part := range gc.Parts {
		resp, err := model.CountTokens(ctx, part)
		if err != nil {
			return sum, err
		}
		sum += int(resp.TotalTokens)
	}
	return sum, nil
}

// SubmitToClient sends all added parts to the specified Vertex AI model for processing,
// returning the model's response. It supports temperature configuration and response trimming.
func (gc *GeminiClient) SubmitToClient(ctx context.Context) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()
	// Configure the model.
	model := gc.Client.GenerativeModel(gc.ModelName)
	model.SetTemperature(gc.Temperature)
	// Pass in the parts and generate a response.
	res, err := model.GenerateContent(ctx, gc.Parts...)
	if err != nil {
		return "", fmt.Errorf("unable to generate contents: %v", err)
	}
	// Examine the response defensively.
	if res == nil || len(res.Candidates) == 0 || res.Candidates[0] == nil ||
		res.Candidates[0].Content == nil || res.Candidates[0].Content.Parts == nil ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from model")
	}
	// Return the result as a string.
	result = fmt.Sprintf("%s\n", res.Candidates[0].Content.Parts[0])
	if gc.Trim {
		return strings.TrimSpace(result), nil
	}
	return result, nil
}

// Submit sends all added parts to the specified Vertex AI model for processing,
// returning the model's response. It supports temperature configuration and response trimming.
// This function creates a temporary client and is not meant to be used within Google Cloud (use SubmitToClient instead).
func (gc *GeminiClient) Submit() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gc.Timeout)
	defer cancel()
	return gc.SubmitToClient(ctx)
}

// CountTokens creates a new client and then counts the tokens in the current multimodal prompt.
func (gc *GeminiClient) CountTokens() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gc.Timeout)
	defer cancel()
	return gc.CountTokensWithContext(ctx)
}

// CountTextTokens tries to count the number of tokens in the given prompt, using the Vertex AI API.
func (gc *GeminiClient) CountTextTokens(prompt string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gc.Timeout)
	defer cancel()
	return gc.CountTextTokensWithClient(ctx, gc.Client, prompt)
}

// CountTextTokensWithClient will count the tokens in the given text
func (gc *GeminiClient) CountTextTokensWithClient(ctx context.Context, client *genai.Client, text string) (int, error) {
	model := client.GenerativeModel(gc.ModelName)
	resp, err := model.CountTokens(ctx, genai.Text(text))
	if err != nil {
		return 0, err
	}
	return int(resp.TotalTokens), nil
}
