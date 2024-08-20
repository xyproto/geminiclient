package simplegemini_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xyproto/env/v2"
	"github.com/xyproto/simplegemini"
)

func TestCustomFunction(t *testing.T) {
	projectID := env.StrAlt("GCP_PROJECT", "PROJECT_ID", "")
	if projectID == "" {
		t.Fatal(simplegemini.ErrGoogleCloudProjectID)
	}

	sf := simplegemini.MustNew()

	// Define a custom function for getting weather
	getWeatherRightNow := func(location string) string {
		fmt.Println("getWeatherRightNow was called")
		switch location {
		case "NY":
			return "It's sunny in New York."
		case "London":
			return "It's rainy in London."
		default:
			return "Weather data not available."
		}
	}

	// Add the function as a tool
	err := sf.AddFunctionTool("get_weather_right_now", "Get the current weather for a specific location", getWeatherRightNow)
	if err != nil {
		t.Fatalf("Failed to add function tool: %v", err)
	}

	// Query Gemini with a prompt that requires using a tool
	result, err := sf.QueryGemini("What is the weather in NY?", nil, nil)
	if err != nil {
		t.Fatalf("Failed to query Gemini: %v", err)
	}

	// Check the response
	if !strings.Contains(result, "It's sunny in New York.") {
		t.Errorf("Expected 'It's sunny in New York.', but got: %v", result)
	}

	fmt.Println("AI Response:", result)
}

func TestReverseStringFunction(t *testing.T) {
	projectID := env.StrAlt("GCP_PROJECT", "PROJECT_ID", "")
	if projectID == "" {
		t.Fatal(simplegemini.ErrGoogleCloudProjectID)
	}

	sf := simplegemini.MustNew()

	// Define a custom function for reversing a string
	reverseString := func(input string) string {
		fmt.Println("reverseString was called")
		runes := []rune(input)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}

	// Add the function as a tool
	err := sf.AddFunctionTool("reverse_string", "Reverse the given string", reverseString)
	if err != nil {
		t.Fatalf("Failed to add function tool: %v", err)
	}

	// Query Gemini with a prompt that requires using the reverse_string function
	result, err := sf.QueryGemini("Reverse the string 'hello'", nil, nil)
	if err != nil {
		t.Fatalf("Failed to query Gemini: %v", err)
	}

	// Check the response
	expected := "olleh"
	if !strings.Contains(result, expected) {
		t.Errorf("Expected '%s', but got: %v", expected, result)
	}

	fmt.Println("Gemini:", result)
}
