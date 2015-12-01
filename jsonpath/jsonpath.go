package jsonpath

import (
	"fmt"
	"reflect"
	"strconv"
)

// JSONPath is an object that can take in a json object, traverse according
// to the rules in the traverser, and return the resulting json. Underlying
// implementation is a singly linked list of Traverser objects
type JSONPath struct {
	head *Traverser
	tail *Traverser
}

// AddTraverser appends a new traverser to the JSONPath's linked list
// of Traverser objects
func (j *JSONPath) AddTraverser(traverser *Traverser) {
	// unitialized JSONPath
	if j.head == nil {
		j.head = traverser
	} else {
		j.tail.child = traverser
	}
	j.tail = traverser
}

// TraverseJSON takes in a json object and returns the subobject specified
// by the JSONPath
func (j *JSONPath) TraverseJSON(json interface{}) (interface{}, error) {
	for {
		if j.head == nil {
			return json, nil
		}
		var err error
		json, err = j.head.Traverse(json)
		if err != nil {
			return nil, err
		}
		j.head = j.head.child
	}
}

// TraverseFunc is a function that takes in json and returns json.
// Should traverse through the input json
type TraverseFunc func(interface{}) (interface{}, error)

// Traverser is really a linked list wrapper over Traverse() functions.
// Traverse() advances through an input json object and returns the result
type Traverser struct {
	child    *Traverser
	Traverse TraverseFunc
}

// NewTraverser returns a new Traverser object
func NewTraverser(f TraverseFunc) *Traverser {
	return &Traverser{Traverse: f}
}

// Key (placeholder) takes in a key name and provides a function to get that key's object
func Key(key string) TraverseFunc {
	return func(json interface{}) (interface{}, error) {
		if jsonMap, ok := json.(map[string]interface{}); ok {
			return getValue(jsonMap, key)
		} else if jsonArray, ok := json.([]interface{}); ok {
			var result []interface{}
			for _, obj := range jsonArray {
				if obj, ok := obj.(map[string]interface{}); ok {
					v, err := getValue(obj, key)
					if err != nil {
						return nil, err
					}
					result = append(result, v)
				} else {
					return nil, fmt.Errorf("expected map[string]interface{} for json array item: %v, got: %v", obj, reflect.TypeOf(obj))
				}
			}
			return result, nil
		}
		return nil, fmt.Errorf("expected map[string]interface{} for json: %v, got: %v", json, reflect.TypeOf(json))
	}
}

// helper function for Key TraverseFunc
func getValue(json map[string]interface{}, key string) (interface{}, error) {
	if result, ok := json[key]; ok {
		return result, nil
	}
	return nil, fmt.Errorf("did not find key '%v' in body: %v", key, json)
}

// Star (placeholder) just returns all the shit at the current json level because it's a WILDCARD!!
func Star() TraverseFunc {
	return func(json interface{}) (interface{}, error) {
		if jsonMap, ok := json.(map[string]interface{}); ok {
			var responseArr []interface{}
			for _, obj := range jsonMap {
				responseArr = append(responseArr, obj)
			}
			return responseArr, nil
		} else if jsonArray, ok := json.([]interface{}); ok {
			return jsonArray, nil
		}
		return nil, fmt.Errorf("uh...nothing to see here")
	}
}

// IndexKey indexes the json (should be an array)
func IndexKey(query string) TraverseFunc {
	return func(json interface{}) (interface{}, error) {
		if jsonArray, ok := json.([]interface{}); ok {
			i, err := strconv.Atoi(query)
			if err != nil {
				return nil, err
			}
			// special case for my use case
			if len(jsonArray) == 0 {
				return nil, nil
			}

			if i >= len(jsonArray) {
				return nil, fmt.Errorf("index array out of range. actual len = %v", len(jsonArray))
			}
			return jsonArray[i], nil
		}
		return nil, fmt.Errorf("cannot index json object: %v", json)
	}
}
