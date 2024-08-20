package simplegemini_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/xyproto/env/v2"
	"github.com/xyproto/simplegemini"
)

var projectID string

func TestMain(m *testing.M) {
	// Check the project ID once before running any tests
	projectID = env.StrAlt("GCP_PROJECT", "PROJECT_ID", "")
	if projectID == "" {
		fmt.Println(simplegemini.ErrGoogleCloudProjectID)
		os.Exit(1)
	}

	// Run the tests
	os.Exit(m.Run())
}

func TestCustomFunction(t *testing.T) {
	gc := simplegemini.MustNew()

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
	err := gc.AddFunctionTool("get_weather_right_now", "Get the current weather for a specific location", getWeatherRightNow)
	if err != nil {
		t.Fatalf("Failed to add function tool: %v", err)
	}

	// Query Gemini with a prompt that requires using a tool
	result, err := gc.Query("What is the weather in NY?")
	if err != nil {
		t.Fatalf("Failed to query Gemini: %v", err)
	}

	// Check the response
	if !strings.Contains(result, "sunny") {
		t.Errorf("Expected 'sunny' to be in the response, but got: %v", result)
	}

	fmt.Println("AI Response:", result)
}

func TestReverseStringFunction(t *testing.T) {
	gc := simplegemini.MustNew()

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
	err := gc.AddFunctionTool("reverse_string", "Reverse the given string", reverseString)
	if err != nil {
		t.Fatalf("Failed to add function tool: %v", err)
	}

	// Query Gemini with a prompt that requires using the reverse_string function
	result, err := gc.Query("Reverse the string 'hello'. Reply with a single word.")
	if err != nil {
		t.Fatalf("Failed to query Gemini: %v", err)
	}

	// Check the response
	expected := "olleh"
	if !strings.Contains(result, expected) {
		t.Errorf("Expected '%s' to be in the response, but got: %v", expected, result)
	}

	fmt.Println("Gemini:", result)
}

func TestNoFunctionsRegistered(t *testing.T) {
	gc := simplegemini.MustNew()

	// Query Gemini with a prompt without any registered functions
	result, err := gc.Query("What is the capital of France? Reply with a single word.")
	if err != nil {
		t.Fatalf("Failed to query Gemini: %v", err)
	}

	if !strings.Contains(result, "Paris") {
		t.Errorf("Expected 'Paris' to be in the response, but got: %v", result)
	}
}

func TestInvalidFunctionRegistration(t *testing.T) {
	gc := simplegemini.MustNew()

	// Attempt to register an invalid function (non-function type)
	err := gc.AddFunctionTool("invalid_tool", "This should fail", "not_a_function")
	if err == nil {
		t.Fatal("Expected an error when registering a non-function, but got none")
	}

	expectedErr := "provided argument is not a function"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error '%s', but got: %v", expectedErr, err)
	}
}

func TestEmptyPrompt(t *testing.T) {
	gc := simplegemini.MustNew()

	// Query Gemini with an empty prompt
	_, err := gc.Query("")
	if err == nil { // success
		t.Fatal("Expected an error when passing in an empty prompt.")
	}
}

func TestAddImageInvalidPath(t *testing.T) {
	gc := simplegemini.MustNew()

	// Attempt to add an image from an invalid path
	err := gc.AddImage("/non/existent/path.png")
	if err == nil {
		t.Fatal("Expected an error when adding an image from an invalid path, but got none")
	}

	expectedErr := "no such file or directory"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error '%s', but got: %v", expectedErr, err)
	}
}

func TestAddURLInvalid(t *testing.T) {
	gc := simplegemini.MustNew()

	// Attempt to add a URL that does not exist
	err := gc.AddURL("http://invalid.url/nonexistent.png")
	if err == nil {
		t.Fatal("Expected an error when adding an invalid URL, but got none")
	}

	expectedErr := "failed to download the file"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error '%s', but got: %v", expectedErr, err)
	}
}
