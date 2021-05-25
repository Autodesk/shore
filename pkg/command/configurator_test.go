package command_test

import (
	"path"
	"testing"

	integ "github.com/Autodeskshore/integration_tests"
	"github.com/Autodeskshore/pkg/command"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var testPath string = "/test"

func TestFailsWithMissingConfig(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Test
		_, err := command.GetConfigFileOrFlag(deps, "render", "values")
		expectedErrMessage := "Config File \"render\" Not Found in \"[/test]\""
		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, expectedErrMessage, err.Error())

	})
}

func TestReadConfigFile(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"a":"a"}`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "render.json"), []byte(renderConfig), 0644)
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", "values")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestViperFlag(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"b":"b"}`
		viper.Set("values", renderConfig)
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", "values")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestViperFileFlag(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"c":"c"}`
		path := path.Join(testPath, "render2.json")
		viper.Set("values", path)

		afero.WriteFile(deps.Project.FS, path, []byte(renderConfig), 0644)
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", "values")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}
