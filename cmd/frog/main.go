package main

import (
	"fmt"
	"log"

	"github.com/xyproto/simplegemini"
	"github.com/xyproto/wordwrap"
)

func main() {
	const (
		multiModalModelName = "gemini-1.0-pro-vision" // "gemini-1.5-pro" also works, if only text is sent
		temperature         = 0.4
		descriptionPrompt   = "Describe what is common for these two images."
	)

	gc, err := simplegemini.NewMultiModal(multiModalModelName, temperature)
	if err != nil {
		log.Fatalf("Could not initialize the Gemini client with the %s model: %v\n", multiModalModelName, err)
	}

	// Build a prompt
	if err := gc.AddImage("frog.png"); err != nil {
		log.Fatalf("Could not add frog.png: %v\n", err)
	}
	gc.AddURI("gs://generativeai-downloads/images/scones.jpg")
	gc.AddText(descriptionPrompt)

	// Count the tokens that are about to be sent
	tokenCount, err := gc.CountTokens()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Sending %d tokens.\n\n", tokenCount)

	// Submit the images and the text prompt
	response, err := gc.Submit()
	if err != nil {
		log.Fatalln(err)
	}

	// Format and print out the response
	if lines, err := wordwrap.WordWrap(response, 79); err == nil { // success
		for _, line := range lines {
			fmt.Println(line)
		}
		return
	}

	fmt.Println(response)
}
