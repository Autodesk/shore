package config

import (
	"path"
	"testing"

	"github.com/Autodeskshore/pkg/project"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFailsWithMissingConfigFile(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Test
		_, err := GetFileConfig(proj, "render")
		// Assert
		assert.Error(t, err)
	})
}

func TestReadConfigFileJSON(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"a":"a"}`
		afero.WriteFile(proj.FS, path.Join(testPath, "render.json"), []byte(renderConfig), 0644)
		// Test
		values, err := GetFileConfig(proj, "render")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestReadConfigFileYml(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"a":"a"}`
		afero.WriteFile(proj.FS, path.Join(testPath, "render.yml"), []byte(renderConfig), 0644)
		// Test
		values, err := GetFileConfig(proj, "render")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestReadConfigFileYaml(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"a":"a"}`
		afero.WriteFile(proj.FS, path.Join(testPath, "render.yaml"), []byte(renderConfig), 0644)
		// Test
		values, err := GetFileConfig(proj, "render")
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}
