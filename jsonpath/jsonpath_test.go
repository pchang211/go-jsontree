package jsonpath

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getJSONBody(raw string) interface{} {
	var body interface{}
	json.Unmarshal([]byte(raw), &body)
	return body
}

func TestSimpleSelect(t *testing.T) {
	path, err := Parse("$.foo")
	assert.NoError(t, err)
	body := getJSONBody(`{"foo":"bar"}`)
	result, err := path.TraverseJSON(body)
	assert.NoError(t, err)
	assert.Equal(t, "bar", result)
}

func TestIndexArray(t *testing.T) {
	path, err := Parse("$.foo[0]")
	assert.NoError(t, err)
	body := getJSONBody(`{"foo":["bar","baz"]}`)
	result, err := path.TraverseJSON(body)
	assert.NoError(t, err)
	assert.Equal(t, "bar", result)

	outOfRangePath, _ := Parse("$.foo[10]")
	_, err = outOfRangePath.TraverseJSON(body)
	assert.Error(t, err)
}

func TestWildcardSelect(t *testing.T) {
	path, err := Parse("$.foo.*")
	assert.NoError(t, err)
	body := getJSONBody(`{"foo":{"one":"bar","two":"baz"}}`)
	result, err := path.TraverseJSON(body)
	assert.NoError(t, err)
	resultArr, ok := result.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(resultArr))
}
