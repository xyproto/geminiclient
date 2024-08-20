package main

import (
	"fmt"

	"github.com/xyproto/simplegemini"
)

func main() {
	fmt.Println(simplegemini.MustAsk("Write a haiku about cows.", 0.4))
}
