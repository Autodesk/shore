package config

import (
	"os"
	"path"
	"testing"

	"github.com/Autodeskshore/pkg/project"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoadShoreConfigFromFile(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
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
		afero.WriteFile(proj.FS, path.Join(testPath, "shore.json"), []byte(shoreConfigString), os.ModePerm)

		// Test
		shoreConfig, err := LoadShoreConfig(proj)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, "test1test2test3", shoreConfig.Profiles[`default`].(map[string]interface{})[`application`])
	})
}

func TestLoadShoreConfigFromMalformedFile(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
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
		afero.WriteFile(proj.FS, path.Join(testPath, "shore.json"), []byte(shoreConfigString), os.ModePerm)

		// Test
		shoreConfig, err := LoadShoreConfig(proj)

		// Assert
		assert.NotNil(t, err)
		assert.Nil(t, shoreConfig.Executor)
		assert.Nil(t, shoreConfig.Renderer)
		assert.Nil(t, shoreConfig.Profiles)
	})
}

func TestLoadShoreConfigFromIncorrectFile(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		shoreConfigString := `{
			"executor": "test",
			"render": "test",
			"profiles": "test"
		}`
		afero.WriteFile(proj.FS, path.Join(testPath, "shore.json"), []byte(shoreConfigString), os.ModePerm)

		// Test
		shoreConfig, err := LoadShoreConfig(proj)

		// Assert
		assert.NotNil(t, err)
		assert.Nil(t, shoreConfig.Executor)
		assert.Nil(t, shoreConfig.Renderer)
		assert.Nil(t, shoreConfig.Profiles)
	})
}

func TestLoadShoreConfigDefault(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		afero.WriteFile(proj.FS, path.Join(testPath, "render.yaml"), []byte(""), os.ModePerm)
		afero.WriteFile(proj.FS, path.Join(testPath, "exec.yml"), []byte(""), os.ModePerm)
		afero.WriteFile(proj.FS, path.Join(testPath, "E2E.json"), []byte(""), os.ModePerm)

		// Test
		shoreConfig, err := LoadShoreConfig(proj)
		shoreConfigExists, configErr := afero.Exists(proj.FS, path.Join(testPath, "shore.yml"))

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, path.Join(testPath, "render.yaml"), shoreConfig.Profiles[`default`].(map[string]interface{})[`render`])
		assert.Equal(t, path.Join(testPath, "exec.yml"), shoreConfig.Profiles[`default`].(map[string]interface{})[`exec`])
		assert.Equal(t, path.Join(testPath, "E2E.json"), shoreConfig.Profiles[`default`].(map[string]interface{})[`e2e`])
		assert.Nil(t, configErr)
		assert.True(t, shoreConfigExists)
	})
}

func TestLoadShoreConfigMissingDefaultConfigs(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Test
		_, errRender := LoadShoreConfig(proj)
		afero.WriteFile(proj.FS, path.Join(testPath, "render.yaml"), []byte(""), os.ModePerm)
		_, errExec := LoadShoreConfig(proj)
		afero.WriteFile(proj.FS, path.Join(testPath, "exec.yml"), []byte(""), os.ModePerm)
		_, errE2E := LoadShoreConfig(proj)

		// Assert
		assert.NotNil(t, errRender)
		assert.EqualError(t, errRender, `unable to find a render config in the project`)
		assert.NotNil(t, errExec)
		assert.EqualError(t, errExec, `unable to find a exec config in the project`)
		assert.NotNil(t, errE2E)
		assert.EqualError(t, errE2E, `unable to find a E2E config in the project`)
	})
}

func TestLoadShoreConfigBadProjectPath(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		os.Setenv("SHORE_PROJECT_PATH", "/tmp/test")

		// Test
		_, err := LoadShoreConfig(proj)

		// Assert
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), `file does not exist`)
	})
}

func TestLoadSpecificConfigFile(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"a":"a"}`
		afero.WriteFile(proj.FS, path.Join(testPath, "render.yaml"), []byte(renderConfig), 0644)
		// Test
		values, err := LoadConfig(proj, "", "render")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestLoadSpecificConfigFlag(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"b":"b"}`
		// Test
		values, err := LoadConfig(proj, renderConfig, "render")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestLoadSpecificConfigFlagOnly(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"a":"a"}`
		renderConfigFileData := `{"b":"b"}`
		afero.WriteFile(proj.FS, path.Join(testPath, "render.yaml"), []byte(renderConfigFileData), 0644)
		// Test
		values, err := LoadConfig(proj, renderConfig, "render")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestLoadSpecificConfigErr(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Test
		_, err := LoadConfig(proj, "", "render")
		// Assert
		assert.Error(t, err)
	})
}
