package jsonnet_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/Autodeskshore/pkg/renderer/jsonnet"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const testPath string = "/tmp/test"

func SetupRenderWithArgs(extension, codeFile, args string) afero.Fs {
	localFs := afero.NewMemMapFs()
	localFs.Mkdir(testPath, os.ModePerm)

	afero.WriteFile(localFs, filepath.Join(testPath, jsonnet.MainFileName), []byte(codeFile), os.ModePerm)

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
	res, renderErr := jsonnet.NewRenderer(fs, logrus.New()).Render(testPath, "")

	// Assert
	assert.Nil(t, renderErr)
	assert.Contains(t, res, `Magic!!!`)
}

func TestFileImporterSuccess(t *testing.T) {
	// Setup Jsonnet Bundler File
	jbFile := `
{
	"version": 1,
	"dependencies": [
	{
		"source": {
		"git": {
			"remote": "https://github.com/org-1/sharedlib1.git",
			"subdir": ""
		}
		},
		"version": "master"
	},
	{
		"source": {
		"git": {
			"remote": "https://github.com/org-2/sharedLib2.git",
			"subdir": ""
		}
		},
		"version": "master"
	},
	{
		"source": {
		"git": {
			"remote": "https://github.com/org-1/sharedLib3.git",
			"subdir": ""
		}
		},
		"version": "master"
	}
	],
	"legacyImports": false
}`

	spec, _ := jsonnetfile.Unmarshal([]byte(jbFile))

	// Test
	importer, err := jsonnet.GetFileImporter(testPath, spec)

	// Assert
	basePath := filepath.Join(testPath, jsonnet.ShareLibsPath, "github.com")

	value := []string{filepath.Join(basePath, "org-1"), filepath.Join(basePath, "org-2")}
	sort.Strings(value)
	sort.Strings(importer.JPaths)

	assert.Len(t, importer.JPaths, 2)
	assert.Nil(t, err)
	assert.Equal(t, value, importer.JPaths)
}

func TestFileImporterLegacySuccess(t *testing.T) {
	// Setup Jsonnet Bundler File
	jbFile := `
{
	"version": 1,
	"dependencies": [
	{
		"source": {
		"git": {
			"remote": "https://github.com/Autodesksharedlib1.git",
			"subdir": ""
		}
		},
		"version": "master"
	},
	{
		"source": {
		"git": {
			"remote": "https://github.com/AutodesksharedLibPath1.git",
			"subdir": ""
		}
		},
		"version": "master"
	}
	],
	"legacyImports": true
}`

	spec, _ := jsonnetfile.Unmarshal([]byte(jbFile))

	// Test
	importer, err := jsonnet.GetFileImporter(testPath, spec)

	// Assert
	value := []string{filepath.Join(testPath, jsonnet.ShareLibsPath)}
	assert.Len(t, importer.JPaths, 1)
	assert.Nil(t, err)
	assert.Equal(t, value, importer.JPaths)
}
