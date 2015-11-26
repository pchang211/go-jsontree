package jsonpath

import (
	"fmt"
	"reflect"
)

// TKey (placeholder) takes in a key name and provides a function to get that key's object
func TKey(key string) TraverseFunc {
	return func(json interface{}) (interface{}, error) {
		if json, ok := json.(map[string]interface{}); ok {
			if result, ok := json[key]; ok {
				return result, nil
			}
			return nil, fmt.Errorf("did not find key '%v' in body: %v", key, json)
		}
		return nil, fmt.Errorf("expected map[string]interface{} for json: %v, got: %v", json, reflect.TypeOf(json))
	}
}
