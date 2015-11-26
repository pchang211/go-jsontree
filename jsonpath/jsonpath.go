package jsonpath

import (
	"fmt"
	"reflect"
	"strconv"
)

// TKey (placeholder) takes in a key name and provides a function to get that key's object
func TKey(key string) TraverseFunc {
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

// TStar (placeholder) just returns all the shit at the current json level because it's a WILDCARD!!
func TStar() TraverseFunc {
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
			if i >= len(jsonArray) {
				return nil, fmt.Errorf("index array out of range. actual len = %v", len(jsonArray))
			}
			return jsonArray[i], nil
		}
		return nil, fmt.Errorf("cannot index json object: %v", json)
	}
}
