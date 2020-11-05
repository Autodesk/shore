package project

import (
	"fmt"
	"os"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func SetupArgsFile(extension, args string) afero.Fs {
	localFs := afero.NewMemMapFs()
	localFs.Mkdir("/tmp/test", os.ModePerm)

	argsFile := fmt.Sprintf("/tmp/test/render.%s", extension)
	afero.WriteFile(localFs, argsFile, []byte(args), os.ModePerm)
	return localFs
}

func TestReadArgsFileYml(t *testing.T) {
	// Given
	argsFile := `
a: test1
b: test2
`
	localFs := SetupArgsFile("yml", argsFile)

	// Test
	project := NewShoreProject(localFs, "/tmp/test/")
	args, err := project.GetRenderArgs()

	// Assert
	assert.Nil(t, err)
	assert.Contains(t, args, `"a":"test1"`)
	assert.Contains(t, args, `"b":"test2"`)
}

func TestReadArgsFileJSON(t *testing.T) {
	// Given
	argsFile := `{"a": "test1", "b": "test2"}`
	localFs := SetupArgsFile("yml", argsFile)

	// Test
	project := NewShoreProject(localFs, "/tmp/test/")
	args, err := project.GetRenderArgs()

	// Assert
	assert.Nil(t, err)
	assert.Contains(t, args, `"a":"test1"`)
	assert.Contains(t, args, `"b":"test2"`)
}

func TestReadArgsFileWithNestedValues(t *testing.T) {
	// Following a bug where nested values weren't passed correctly and would cause a panic.
	// This test validates this bug never returns.

	// Given
	argsFile := `{
	"a": "test1",
	"b": "test2",
	"c": [
		"d": {
			"e": "f",
			"g": [
				1,
				2,
				3,
				[4, "5", 6],
				{"7":8},
				[9, 10],
			]
		}
	]
}`
	localFs := SetupArgsFile("json", argsFile)

	// Test
	project := NewShoreProject(localFs, "/tmp/test/")
	args, err := project.GetRenderArgs()

	// Assert
	argsMap := make(map[string]interface{})

	if unMarshalErr := jsoniter.Unmarshal([]byte(args), &argsMap); unMarshalErr != nil {
		assert.Nil(t, unMarshalErr)
	}

	assert.Nil(t, err)
	assert.Contains(t, args, "test1")
	assert.Contains(t, args, "test2")
	assert.Contains(t, args, "test2")
	assert.Contains(t, args, `{"7":8}`)
	assert.Contains(t, args, `[9,10]`)
}

func TestNoArgsFileReutnrsEmptyAndNil(t *testing.T) {
	// Given
	localFs := afero.NewMemMapFs()

	// Test
	project := NewShoreProject(localFs, "/tmp/test/")
	args, err := project.GetRenderArgs()

	// Assert
	assert.Equal(t, "", args)
	assert.Equal(t, true, os.IsNotExist(err))
}
