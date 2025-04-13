package k8srender_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/Autodesk/shore/pkg/renderer/jsonnet"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const testPath string = "/tmp/test"

func SetupRenderWithArgs(extension, codeFile, args string) afero.Fs {
	localFs := afero.NewMemMapFs()
	localFs.Mkdir(testPath, os.ModePerm)

	afero.WriteFile(localFs, filepath.Join(testPath, jsonnet.RenderFiles[renderer.MainFileName]), []byte(codeFile), os.ModePerm)

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
	res, renderErr := jsonnet.NewRenderer(fs, logrus.New()).Render(testPath, "", renderer.MainFileName)

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
	importer := jsonnet.NewImporter(&afero.MemMapFs{}, testPath, spec)

	// Assert
	basePath := filepath.Join(testPath, jsonnet.ShareLibsPath, "github.com")

	value := []string{testPath, filepath.Join(basePath, "org-1"), filepath.Join(basePath, "org-2")}
	sort.Strings(value)
	sort.Strings(importer.JPaths)

	assert.Len(t, importer.JPaths, 3)
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
			"remote": "https://github.com/example/sharedlib1.git",
			"subdir": ""
		}
		},
		"version": "master"
	},
	{
		"source": {
		"git": {
			"remote": "https://github.com/example/sharedLibPath1.git",
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
	importer := jsonnet.NewImporter(&afero.MemMapFs{}, testPath, spec)

	// Assert
	value := []string{testPath, filepath.Join(testPath, jsonnet.ShareLibsPath)}
	assert.Len(t, importer.JPaths, 2)
	assert.Equal(t, value, importer.JPaths)
}

func TestRenderOutputIsYAML(t *testing.T) {
	// Given
	codeFile := `
{
  apiVersion: "v1",
  kind: "ConfigMap",
  metadata: {
    name: "example-config",
    namespace: "default",
  },
  data: {
    key1: "value1",
    key2: "value2",
  },
}
`
	fs := SetupRenderWithArgs("jsonnet", codeFile, "")

	// Test
	output, err := jsonnet.NewRenderer(fs, logrus.New()).Render(testPath, "", renderer.MainFileName)

	// Assert
	assert.NoError(t, err, "Render should not return an error")

	// Validate that the output is valid YAML
	var yamlData map[string]interface{}
	err = yaml.Unmarshal([]byte(output), &yamlData)
	assert.NoError(t, err, "Output should be valid YAML")

	// Additional assertions to verify the content
	assert.Equal(t, "v1", yamlData["apiVersion"])
	assert.Equal(t, "ConfigMap", yamlData["kind"])
	assert.Equal(t, "example-config", yamlData["metadata"].(map[interface{}]interface{})["name"])
	assert.Equal(t, "default", yamlData["metadata"].(map[interface{}]interface{})["namespace"])
	assert.Equal(t, "value1", yamlData["data"].(map[interface{}]interface{})["key1"])
	assert.Equal(t, "value2", yamlData["data"].(map[interface{}]interface{})["key2"])
}
