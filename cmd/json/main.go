package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xyproto/geminiclient"
)

func main() {
	const (
		prompt      = `What color is the sky? Answer with a JSON struct where the only key is "color" and the value is a lowercase string.`
		modelName   = "gemini-1.5-pro"
		temperature = 0.0
		timeout     = 10 * time.Second
	)

	gc, err := geminiclient.NewWithTimeout(modelName, temperature, timeout)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(prompt)

	result, err := gc.Query(prompt)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(result)
}
