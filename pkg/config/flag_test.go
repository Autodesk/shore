package config

import (
	"path"
	"testing"

	"github.com/Autodeskshore/pkg/project"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFlagWithJson(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"b":"b"}`
		// Test
		values, err := GetFlagConfig(proj, renderConfig)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithFilepath(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"c":"c"}`
		fileName := "render2.json"
		path := path.Join(testPath, fileName)
		afero.WriteFile(proj.FS, path, []byte(renderConfig), 0644)

		// Test
		values, err := GetFlagConfig(proj, fileName)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithFilepathAbs(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"c":"c"}`
		path := path.Join(testPath, "render2.json")
		afero.WriteFile(proj.FS, path, []byte(renderConfig), 0644)

		// Test
		values, err := GetFlagConfig(proj, path)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithUpperCase(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"B":"b"}`
		// Test
		values, err := GetFlagConfig(proj, renderConfig)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithDottedKeys(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"B.C.E":"b"}`
		// Test
		values, err := GetFlagConfig(proj, renderConfig)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithMessyData(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `{"B.C.E":"b", "ABC123": "a", "1": 1, "2": 2}`
		// Test
		values, err := GetFlagConfig(proj, renderConfig)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, renderConfig, string(values))
	})
}

func TestFlagWithBadData(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		renderConfig := `yolo swaggins bilbo baggins`
		// Test
		_, err := GetFlagConfig(proj, renderConfig)
		// Assert
		assert.Error(t, err)
	})
}

func TestFlagWithBadJson(t *testing.T) {
	SetupTest(t, func(t *testing.T, proj *project.Project) {
		// Given
		// Missing a `,` after the `"b"`
		renderConfig := `{"B.C.E":"b" "ABC123": "a", "1": 1, "2": 2}`
		// Test
		_, err := GetFlagConfig(proj, renderConfig)
		// Assert
		assert.Error(t, err)
	})
}
