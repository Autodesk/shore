package jsonnet

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	// Given
	fileName := "test.jsonnet"
	code := map[string]string{
		"Is":      "This",
		"Working": "false",
	}

	codeBytes, err := json.Marshal(code)

	// Test
	renderer := NewRenderer()
	res, err := renderer.Render(fileName, string(codeBytes))

	var renderedCode map[string]string
	json.Unmarshal([]byte(res), &renderedCode)

	// Assert
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(code, renderedCode) { // If both maps don't match, we should throw an error.
		err := `
Expected code to match rendered code,

Code: %+v,
RenderedCode: %+v
`
		t.Errorf(err, code, renderedCode)
	}
}

func TestSharedLibraryLoad(t *testing.T) {
	// Given
	fileName := "test.jsonnet"
	code := `
local pipeline = "sponnet/pipeline.libsonnet";

pipeline.Pipeline {}
`

	codeBytes, err := json.Marshal(code)

	// Test
	path, err := os.Getwd()
	shreadLibPath := filepath.Join(path, "..", "..", "..")
	renderer := NewRenderer(shreadLibPath)
	res, err := renderer.Render(fileName, string(codeBytes))

	var renderedCode map[string]string
	json.Unmarshal([]byte(res), &renderedCode)

	// Assert
	if err != nil {
		t.Error(err)
	}
}
