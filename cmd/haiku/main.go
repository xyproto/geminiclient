package main

import (
	"fmt"

	"github.com/xyproto/geminiclient"
)

func main() {
	fmt.Println(geminiclient.MustAsk("Write a haiku about cows.", 0.4))
}
