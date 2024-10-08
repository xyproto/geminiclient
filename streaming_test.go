package geminiclient_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/xyproto/geminiclient"
)

// TestSubmitToClientStreaming tests the SubmitToClientStreaming function.
func TestSubmitToClientStreaming(t *testing.T) {
	// Initialize a new GeminiClient
	gc := geminiclient.MustNewWithTimeout("gemini-1.5-pro", 0.0, 10*time.Second)

	// Create a context
	ctx := context.Background()

	// Define a prompt and add it as text to the GeminiClient
	prompt := "Write a story about a magic backpack, around 50 words. Make sure to include the phrase \"magic backpack\"."
	gc.AddText(prompt)

	// Capture the streamed content
	var streamedContent strings.Builder
	streamCallback := func(part string) {
		fmt.Println("Streaming part:", part) // Debugging output
		streamedContent.WriteString(part)
	}

	// Run the SubmitToClientStreaming function
	_, err := gc.SubmitToClientStreaming(ctx, streamCallback)
	if err != nil {
		t.Fatalf("Streaming failed: %v", err)
	}

	// Check if the streamed content is not empty
	if streamedContent.Len() == 0 {
		t.Fatal("Expected streamed content, but got empty result")
	}

	// Optional: Check if the streamed content contains expected text
	if !strings.Contains(streamedContent.String(), "magic backpack") {
		t.Fatalf("Expected streamed content to contain 'magic backpack', but it didn't. Actual content: %s", streamedContent.String())
	}

	fmt.Println("Streamed content:", streamedContent.String())
}
