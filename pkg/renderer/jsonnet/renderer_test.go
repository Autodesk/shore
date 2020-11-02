package jsonnet

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Autodesk/shore/pkg/fs"
	"github.com/google/go-jsonnet"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewRenderer(t *testing.T) {
	// Given
	localFs := fs.InitFs(fs.MEM)
	localFs.Mkdir("/tmp/test", os.ModePerm)

	code := map[string]string{
		"Is":      "This",
		"Working": "false",
	}

	codeBytes, marshalErr := json.Marshal(code)

	afero.WriteFile(localFs, "/tmp/test/main.pipeline.jsonnet", codeBytes, os.ModePerm)

	// Test
	renderer, rendererErr := NewRenderer("/tmp/test")
	res, renderErr := renderer.Render()

	var renderedCode map[string]string
	json.Unmarshal([]byte(res), &renderedCode)

	// Assert
	assert.Nil(t, marshalErr)
	assert.Nil(t, rendererErr)
	assert.Nil(t, renderErr)
	assert.Equal(t, code, renderedCode, `
Expected code to match rendered code,

Code: %+v,
RenderedCode: %+v
`)
}

type TestImporter struct {
	JPaths []string
}

func (i *TestImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	localFs := fs.GetFs()
	for _, p := range i.JPaths {
		libPath := filepath.Join(p, importedPath)
		if exists, err := afero.Exists(localFs, libPath); err == nil && exists {
			lib, err := afero.ReadFile(localFs, libPath)

			if err != nil {
				return jsonnet.MakeContents(""), "", err
			}

			return jsonnet.MakeContents(string(lib)), libPath, nil
		}
	}

	return jsonnet.MakeContents(""), "", fmt.Errorf("Library %s not found", importedPath)
}

func TestSharedLibraryLoad(t *testing.T) {
	// Given
	localFs := fs.InitFs(fs.MEM)
	projectPath := "/tmp/test"
	vendorPath := filepath.Join(projectPath, "vendor")
	sharedLibPath := filepath.Join(vendorPath, "sharedlib")
	sharedLibDirPath := filepath.Join(vendorPath, "sharedlib", "sponnet")

	localFs.Mkdir(projectPath, os.ModePerm)
	localFs.Mkdir(vendorPath, os.ModePerm)
	localFs.Mkdir(sharedLibPath, os.ModePerm)
	localFs.Mkdir(sharedLibDirPath, os.ModePerm)

	code := `
local pipeline = import "sponnet/pipeline.libsonnet";

pipeline.Pipeline{}
`

	goMod := `
module github.com/Autodesk/sponnet-v1-test

go 1.15

require sharedlib/sponnet v1.0.0
`

	sponnetLib := `
{
	Pipeline: {"This": "Works"},
}
`

	afero.WriteFile(localFs, filepath.Join(projectPath, "main.pipeline.jsonnet"), []byte(code), os.ModePerm)
	afero.WriteFile(localFs, filepath.Join(projectPath, "go.mod"), []byte(goMod), os.ModePerm)
	afero.WriteFile(localFs, filepath.Join(sharedLibDirPath, "pipeline.libsonnet"), []byte(sponnetLib), os.ModePerm)

	renderer := &Jsonnet{VM: jsonnet.MakeVM(), projectPath: projectPath, FS: localFs}
	sharedLibs, sharedLibsErr := getSharedLibs(localFs, projectPath)
	renderer.VM.Importer(&TestImporter{JPaths: sharedLibs})

	// Test
	res, renderErr := renderer.Render()

	var renderedCode map[string]string
	unmarshalErr := json.Unmarshal([]byte(res), &renderedCode)

	// Assert
	assert.Nil(t, sharedLibsErr)
	assert.Nil(t, renderErr)
	assert.Nil(t, unmarshalErr)
	assert.Equal(t, renderedCode, map[string]string{"This": "Works"})
}
