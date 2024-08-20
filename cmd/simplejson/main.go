package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xyproto/simplegemini"
)

func main() {
	const (
		prompt      = `What color is the sky? Answer with a JSON struct where the only key is "color" and the value is a lowercase string.`
		modelName   = "gemini-1.5-pro" // "gemini-1.5-flash"
		temperature = 0.0
		timeout     = 10 * time.Second
	)

	geminiClient, err := simplegemini.NewWithTimeout(modelName, temperature, timeout)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(prompt)

	result, err := geminiClient.Query(prompt)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(result)
}
