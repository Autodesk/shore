package jsonnet

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/modfile"
)

const testPath string = "/tmp/test"

func SetupRenderWithArgs(extension, codeFile, args string) afero.Fs {
	localFs := afero.NewMemMapFs()
	localFs.Mkdir(testPath, os.ModePerm)

	afero.WriteFile(localFs, filepath.Join(testPath, MainFileName), []byte(codeFile), os.ModePerm)

	if args != "" {
		argsFile := filepath.Join(testPath, fmt.Sprintf("render.%s", extension))
		afero.WriteFile(localFs, argsFile, []byte(args), os.ModePerm)
	}

	return localFs
}

func TestNewRenderer(t *testing.T) {
	// Given
	codeFile := `
function(params={})(
	{"This-is": "Magic!!!"}
)
`
	fs := SetupRenderWithArgs("json", codeFile, "")

	// Test
	res, renderErr := NewRenderer(fs, logrus.New()).Render(testPath, "")

	// Assert
	assert.Nil(t, renderErr)
	assert.Contains(t, res, `Magic!!!`)
}

// type TestImporter struct {
// 	JPaths []string
// }

// func (i *TestImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
// 	localFs := fs.GetFs()
// 	for _, p := range i.JPaths {
// 		libPath := filepath.Join(p, importedPath)
// 		if exists, err := afero.Exists(localFs, libPath); err == nil && exists {
// 			lib, err := afero.ReadFile(localFs, libPath)

// 			if err != nil {
// 				return jsonnet.MakeContents(""), "", err
// 			}

// 			return jsonnet.MakeContents(string(lib)), libPath, nil
// 		}
// 	}

// 	return jsonnet.MakeContents(""), "", fmt.Errorf("Library %s not found", importedPath)
// }

func SetupSharedLibs() afero.Fs {
	localFs := afero.NewMemMapFs()
	vendorPath := filepath.Join(testPath, "vendor")
	sharedLibPath := filepath.Join(vendorPath, "sharedlib")
	sharedLibDirPath := filepath.Join(vendorPath, "sharedlib", "sponnet")
	sharedLib2Path := filepath.Join(vendorPath, "sharedlib2")
	sharedLibDir2Path := filepath.Join(vendorPath, "sharedlib2", "sponnet2")

	localFs.Mkdir(testPath, os.ModePerm)
	localFs.Mkdir(vendorPath, os.ModePerm)
	localFs.Mkdir(sharedLibPath, os.ModePerm)
	localFs.Mkdir(sharedLibDirPath, os.ModePerm)
	localFs.Mkdir(sharedLib2Path, os.ModePerm)
	localFs.Mkdir(sharedLibDir2Path, os.ModePerm)

	goMod := `
module github.com/Autodesk/generic-pipeline

go 1.15

require sharedlib/sponnet v1.0.0
require sharedlib2/sponnet2 v1.0.0
`

	sponnetLib := `
{
	Pipeline: {"This": "Works"},
}
`

	afero.WriteFile(localFs, filepath.Join(testPath, "go.mod"), []byte(goMod), os.ModePerm)
	afero.WriteFile(localFs, filepath.Join(sharedLibDirPath, "pipeline.libsonnet"), []byte(sponnetLib), os.ModePerm)
	afero.WriteFile(localFs, filepath.Join(sharedLibDir2Path, "pipeline.libsonnet"), []byte(sponnetLib), os.ModePerm)
	return localFs
}

func TestSharedLibraryLoad(t *testing.T) {
	// Given
	localFs := SetupSharedLibs()

	// Test
	renderer := NewRenderer(localFs, logrus.New())
	libs, err := renderer.getSharedLibs(testPath)

	// Assert
	assert.Nil(t, err)
	assert.Len(t, libs, 2)
}

func TestSharedLibraryLoadNoModFile(t *testing.T) {
	// Given
	localFs := afero.NewMemMapFs()

	// Test
	renderer := NewRenderer(localFs, logrus.New())
	libs, err := renderer.getSharedLibs(testPath)

	// Assert
	assert.Len(t, libs, 0)
	assert.Equal(t, os.IsNotExist(err), true)
}

func TestPipelineFileDoesNotExist(t *testing.T) {
	// Given
	localFs := afero.NewMemMapFs()

	// Test
	renderer := NewRenderer(localFs, logrus.New())
	pipeline, err := renderer.Render(testPath, "")

	// Assert
	assert.Equal(t, pipeline, "")
	assert.Error(t, err)
}

func TestSharedLibraryLoadModFileNotExist(t *testing.T) {
	// Given
	localFs := SetupSharedLibs()
	localFs.Remove(filepath.Join(testPath, "go.mod"))

	// Test
	renderer := NewRenderer(localFs, logrus.New())
	libs, err := renderer.getSharedLibs(testPath)

	// Assert
	assert.Len(t, libs, 0)
	assert.Equal(t, os.IsNotExist(err), true)
}

func TestSharedLibraryLoadModBroken(t *testing.T) {
	// Given
	localFs := SetupSharedLibs()

	// This is a broken goMod.
	goMod := `
module github.com/Autodesk/sponnet-v1-test

go 1.15

required sharedlib/sponnet v1.0.0
`
	// Override the modfile.
	afero.WriteFile(localFs, filepath.Join(testPath, "go.mod"), []byte(goMod), os.ModePerm)

	// Test
	renderer := NewRenderer(localFs, logrus.New())
	libs, err := renderer.getSharedLibs(testPath)

	// Assert
	assert.Len(t, libs, 0)
	assert.IsType(t, modfile.ErrorList{}, err)
}

func TestSharedLibraryMissing(t *testing.T) {
	// Given
	localFs := SetupSharedLibs()

	// Override the modfile.
	vendorPath := filepath.Join(testPath, "vendor")
	sharedLibPath := filepath.Join(vendorPath, "sharedlib")
	localFs.Remove(sharedLibPath)

	// Test
	renderer := NewRenderer(localFs, logrus.New())
	libs, err := renderer.getSharedLibs(testPath)

	// Assert
	assert.IsType(t, SharedLibsErr{}, err)
	// We deleted `sharedlib` but kept `sharelib2`.
	// The method should return the libraries it found and those that it DIDN'T find
	assert.Len(t, libs, 1)

}
