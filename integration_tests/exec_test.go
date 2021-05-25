package integration_tests

import (
	"os"
	"path"
	"testing"

	"github.com/Autodeskshore/pkg/command"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulExecWithConfigFile(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execConfig := `{"application": "First Application", "pipeline": "First Pipeline"}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "exec.json"), []byte(execConfig), os.ModePerm)

		// Test
		execCmd := command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulExecWithFlag(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execConfig := `{"application": "First Application", "pipeline": "First Pipeline"}`

		viper.Set("payload", execConfig)
		command.GetConfigFileOrFlag(deps, "exec", "payload")

		// Test
		execCmd := command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestFailureExecWithConfigFileMissingParameter(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execError := "required args key 'pipeline' missing"
		execConfig := `{"application": "First Application"}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "exec.json"), []byte(execConfig), os.ModePerm)

		// Test
		execCmd := command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}

func TestFailureExecWithFlagMissingParameter(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execError := "required args key 'pipeline' missing"
		execConfig := `{"application": "First Application"}`

		viper.Set("payload", execConfig)
		command.GetConfigFileOrFlag(deps, "exec", "payload")

		// Test
		execCmd := command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}

func TestFailureExecWithConfigFileBadPayload(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execError := "ReadMapCB: expect { or n, but found \x00, error found in #0 byte of ...||..., bigger context ...||..."
		execConfig := `{"application": "First Application",}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "exec.json"), []byte(execConfig), os.ModePerm)

		// Test
		execCmd := command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}

func TestFailureExecWithFlagBadPayload(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execError := "required args key 'pipeline' missing\nrequired args key 'application' missing"
		execConfig := `{"application": "First Application",  "pipeline": "First Pipeline" ,,,,,,}`

		viper.Set("payload", execConfig)
		command.GetConfigFileOrFlag(deps, "exec", "payload")

		// Test
		execCmd := command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}
