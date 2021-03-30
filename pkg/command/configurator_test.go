package command_test

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/Autodeskshore/pkg/backend/spinnaker"
	"github.com/Autodeskshore/pkg/command"
	"github.com/Autodeskshore/pkg/project"
	"github.com/Autodeskshore/pkg/renderer/jsonnet"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var Deps *command.Dependencies
var tmpDir string

func init() {
	// Move the working dir to /tmp since we may hit the local file system by accident.
	tmpDir = "/tmp"
	if runtime.GOOS == "darwin" {
		tmpDir = "/private/tmp"
	}

	os.Chdir(tmpDir)
	fs := afero.NewMemMapFs()
	logger, _ := test.NewNullLogger()

	Deps = &command.Dependencies{
		Project:  project.NewShoreProject(fs, logger),
		Renderer: jsonnet.NewRenderer(fs, logger),
		Backend:  spinnaker.NewClient(logger),
		Logger:   logger,
	}
}

func TestFailsWithMissingConfig(t *testing.T) {
	// Test
	_, err := command.GetConfigFileOrFlag(Deps, "render", "values")
	// Assert
	assert.NotNil(t, err)
}

func TestReadConfigFile(t *testing.T) {
	// Given
	data := `{"a":"a"}`
	afero.WriteFile(Deps.Project.FS, fmt.Sprintf("%v/render.json", tmpDir), []byte(data), 0644)
	// Test
	values, err := command.GetConfigFileOrFlag(Deps, "render", "values")
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, data, string(values))
}

func TestViperFlag(t *testing.T) {
	// Given
	data := `{"b":"b"}`
	viper.Set("values", data)
	// Test
	values, err := command.GetConfigFileOrFlag(Deps, "render", "values")
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, data, string(values))
}

func TestViperFileFlag(t *testing.T) {
	// Given
	data := `{"c":"c"}`
	path := fmt.Sprintf("%v/render2.json", tmpDir)
	viper.Set("values", path)

	afero.WriteFile(Deps.Project.FS, path, []byte(data), 0644)
	// Test
	values, err := command.GetConfigFileOrFlag(Deps, "render", "values")
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, string(values), data)
}
