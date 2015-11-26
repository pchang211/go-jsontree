package main

import (
	"fmt"
	"os"
	"pchang211/go-jsontree/exp/jsonpath"
)

func main() {
	input := os.Args[1]
	_, err := jsonpath.Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v", err)
	}
}
