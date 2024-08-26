package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/xyproto/geminiclient"
)

func main() {
	gc := geminiclient.MustNew()

	// Define a custom function for getting the weather, that Gemini can choose to call
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

	// Add the weather function as a tool
	err := gc.AddFunctionTool("get_weather_right_now", "Get the current weather for a specific location", getWeatherRightNow)
	if err != nil {
		log.Fatalf("Failed to add function tool: %v", err)
	}

	// Query Gemini with a prompt that requires using the custom weather tool
	result, err := gc.Query("What is the weather in NY?")
	if err != nil {
		log.Fatalf("Failed to query Gemini: %v", err)
	}

	// Check and print the weather response
	if !strings.Contains(result, "sunny") {
		log.Fatalf("Expected 'sunny' to be in the response, but got: %v", result)
	}
	fmt.Println("Weather AI Response:", result)

	gc.Clear() // Clear the current prompt parts, tools and functions

	// Define a custom function for reversing a string
	reverseString := func(input string) string {
		fmt.Println("reverseString was called")
		runes := []rune(input)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}

	// Add the string reversal function as a tool
	err = gc.AddFunctionTool("reverse_string", "Reverse the given string", reverseString)
	if err != nil {
		log.Fatalf("Failed to add function tool: %v", err)
	}

	// Query Gemini with a prompt that requires using the string reversal tool
	result, err = gc.Query("Reverse the string 'hello'. Reply with a single word.")
	if err != nil {
		log.Fatalf("Failed to query Gemini: %v", err)
	}

	// Check and print the string reversal response
	expected := "olleh"
	if !strings.Contains(result, expected) {
		log.Fatalf("Expected '%s' to be in the response, but got: %v", expected, result)
	}
	fmt.Println("Reversal AI Response:", result)
}
