package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"pchang211/go-jsontree/jsonpath"
)

func main() {

	jsonString := flag.String("j", "", "raw json string")
	flag.Parse()
	args := flag.Args()

	jsonRawInput := []byte(*jsonString)
	var inputJSON interface{}
	err := json.Unmarshal(jsonRawInput, &inputJSON)
	if err != nil {
		fmt.Printf("error reading json: %v\n", err)
		return
	}
	fmt.Printf("json: %v\n", inputJSON)

	path := args[0]
	jsonpath, err := jsonpath.Parse(path)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	result, err := jsonpath.TraverseJSON(inputJSON)
	if err != nil {
		fmt.Printf("error traversing: %v\n", err)
		return
	}
	fmt.Printf("result: %v\n", result)

}
