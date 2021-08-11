package command_test

import (
	"path"
	"testing"

	integ "github.com/Autodeskshore/integration_tests"
	"github.com/Autodeskshore/pkg/command"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var testPath string = "/test"

func TestFailsWithMissingConfig(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Test
		_, err := command.GetConfigFileOrFlag(deps, "render", "")
		// Assert
		assert.Error(t, err)
	})
}

func TestReadConfigFileJSON(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"a":"a"}`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "render.json"), []byte(renderConfig), 0644)
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", "")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestReadConfigFileYml(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"a":"a"}`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "render.yml"), []byte(renderConfig), 0644)
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", "")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestReadConfigFileYaml(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"a":"a"}`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "render.yaml"), []byte(renderConfig), 0644)
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", "")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithJson(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"b":"b"}`
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", renderConfig)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithFilepath(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"c":"c"}`
		fileName := "render2.json"
		path := path.Join(testPath, fileName)

		afero.WriteFile(deps.Project.FS, path, []byte(renderConfig), 0644)
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", fileName)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithFilepathAbs(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"c":"c"}`
		path := path.Join(testPath, "render2.json")

		afero.WriteFile(deps.Project.FS, path, []byte(renderConfig), 0644)
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", path)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithUpperCase(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"B":"b"}`
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", renderConfig)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithDottedKeys(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"B.C.E":"b"}`
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", renderConfig)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithMessyData(t *testing.T) {
	integ.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"B.C.E":"b", "ABC123": "a", "1": 1, "2": 2}`
		// Test
		values, err := command.GetConfigFileOrFlag(deps, "render", renderConfig)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}
