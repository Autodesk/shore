package config

import (
	"os"
	"path"
	"testing"

	"github.com/Autodeskshore/pkg/project"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoadShoreConfigFromFile(t *testing.T) {
	// Given
	memFs := afero.NewMemMapFs()
	testPath := "/tmp/test"
	memFs.Mkdir(testPath, os.ModePerm)
	os.Setenv("LOCAL", "true")
	os.Setenv("SHORE_PROJECT_PATH", testPath)

	shoreConfigString := `{
		"renderer": {
		  "type": "jsonnet"
		},
		"executor": {
		  "type": "spinnaker",
		  "config": {
			"default": "~/.spin/sb-config",
			"prodSpin": "~/.spin/prod-config"
		  }
		},
		"profiles": {
		  "default": {
			"application": "test1test2test3",
			"pipeline": "simple-pipeline-test",
			"render": "render.yaml",
			"exec": "exec.yaml",
			"e2e": "e2e.yaml"
		  }
		}
	}`
	afero.WriteFile(memFs, path.Join(testPath, "shore.json"), []byte(shoreConfigString), os.ModePerm)

	logger, _ := test.NewNullLogger()
	p := project.NewShoreProject(memFs, logger)

	// Test
	shoreConfig, err := LoadShoreConfig(p)
	os.Unsetenv("SHORE_PROJECT_PATH")
	os.Unsetenv("LOCAL")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "test1test2test3", shoreConfig.Profiles[`default`].(map[string]interface{})[`application`])
}

func TestLoadShoreConfigFromMalformedFile(t *testing.T) {
	// Given
	memFs := afero.NewMemMapFs()
	testPath := "/tmp/test"
	memFs.Mkdir(testPath, os.ModePerm)
	os.Setenv("LOCAL", "true")
	os.Setenv("SHORE_PROJECT_PATH", testPath)

	shoreConfigString := `{
		"renderer": {
		  "type": "jsonnet"
		},
		"executor": {
		  "type": "spinnaker",
		  "config": 
		},
		"profiles": {
		  "default": 
		}
	}`
	afero.WriteFile(memFs, path.Join(testPath, "shore.json"), []byte(shoreConfigString), os.ModePerm)

	logger, _ := test.NewNullLogger()
	p := project.NewShoreProject(memFs, logger)

	// Test
	shoreConfig, err := LoadShoreConfig(p)
	os.Unsetenv("SHORE_PROJECT_PATH")
	os.Unsetenv("LOCAL")

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, shoreConfig.Executor)
	assert.Nil(t, shoreConfig.Renderer)
	assert.Nil(t, shoreConfig.Profiles)

}

func TestLoadShoreConfigFromIncorrectFile(t *testing.T) {
	// Given
	memFs := afero.NewMemMapFs()
	testPath := "/tmp/test"
	memFs.Mkdir(testPath, os.ModePerm)
	os.Setenv("LOCAL", "true")
	os.Setenv("SHORE_PROJECT_PATH", testPath)

	shoreConfigString := `{
		"executor": "test",
		"render": "test",
		"profiles": "test"
	}`
	afero.WriteFile(memFs, path.Join(testPath, "shore.json"), []byte(shoreConfigString), os.ModePerm)

	logger, _ := test.NewNullLogger()
	p := project.NewShoreProject(memFs, logger)

	// Test
	shoreConfig, err := LoadShoreConfig(p)
	os.Unsetenv("SHORE_PROJECT_PATH")
	os.Unsetenv("LOCAL")

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, shoreConfig.Executor)
	assert.Nil(t, shoreConfig.Renderer)
	assert.Nil(t, shoreConfig.Profiles)
}

func TestLoadShoreConfigDefault(t *testing.T) {
	// Given
	memFs := afero.NewMemMapFs()
	testPath := "/tmp/test"
	memFs.Mkdir(testPath, os.ModePerm)
	os.Setenv("LOCAL", "true")
	os.Setenv("SHORE_PROJECT_PATH", testPath)
	afero.WriteFile(memFs, path.Join(testPath, "render.yaml"), []byte(""), os.ModePerm)
	afero.WriteFile(memFs, path.Join(testPath, "exec.yml"), []byte(""), os.ModePerm)
	afero.WriteFile(memFs, path.Join(testPath, "E2E.json"), []byte(""), os.ModePerm)

	logger, _ := test.NewNullLogger()
	p := project.NewShoreProject(memFs, logger)

	// Test
	shoreConfig, err := LoadShoreConfig(p)
	shoreConfigExists, configErr := afero.Exists(memFs, path.Join(testPath, "shore.yml"))
	os.Unsetenv("SHORE_PROJECT_PATH")
	os.Unsetenv("LOCAL")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, path.Join(testPath, "render.yaml"), shoreConfig.Profiles[`default`].(map[string]interface{})[`render`])
	assert.Equal(t, path.Join(testPath, "exec.yml"), shoreConfig.Profiles[`default`].(map[string]interface{})[`exec`])
	assert.Equal(t, path.Join(testPath, "E2E.json"), shoreConfig.Profiles[`default`].(map[string]interface{})[`e2e`])
	assert.Nil(t, configErr)
	assert.True(t, shoreConfigExists)
}

func TestLoadShoreConfigMissingDefaultConfigs(t *testing.T) {
	// Given
	memFs := afero.NewMemMapFs()
	testPath := "/tmp/test"
	memFs.Mkdir(testPath, os.ModePerm)
	os.Setenv("LOCAL", "true")
	os.Setenv("SHORE_PROJECT_PATH", testPath)

	logger, _ := test.NewNullLogger()
	p := project.NewShoreProject(memFs, logger)

	// Test
	_, errRender := LoadShoreConfig(p)
	afero.WriteFile(memFs, path.Join(testPath, "render.yaml"), []byte(""), os.ModePerm)
	_, errExec := LoadShoreConfig(p)
	afero.WriteFile(memFs, path.Join(testPath, "exec.yml"), []byte(""), os.ModePerm)
	_, errE2E := LoadShoreConfig(p)

	os.Unsetenv("SHORE_PROJECT_PATH")
	os.Unsetenv("LOCAL")

	// Assert
	assert.NotNil(t, errRender)
	assert.EqualError(t, errRender, `unable to find a render config in the project`)
	assert.NotNil(t, errExec)
	assert.EqualError(t, errExec, `unable to find a exec config in the project`)
	assert.NotNil(t, errE2E)
	assert.EqualError(t, errE2E, `unable to find a E2E config in the project`)
}

func TestLoadShoreConfigBadProjectPath(t *testing.T) {
	// Given
	memFs := afero.NewMemMapFs()
	os.Setenv("LOCAL", "true")
	os.Setenv("SHORE_PROJECT_PATH", "/tmp/test")

	logger, _ := test.NewNullLogger()
	p := project.NewShoreProject(memFs, logger)

	// Test
	_, err := LoadShoreConfig(p)
	os.Unsetenv("SHORE_PROJECT_PATH")
	os.Unsetenv("LOCAL")

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), `file does not exist`)
}
