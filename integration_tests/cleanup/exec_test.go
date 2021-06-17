package cleanup_test

import (
	"os"
	"path"
	"testing"

	"github.com/Autodeskshore/integration_tests"
	"github.com/Autodeskshore/pkg/cleanup_command"
	"github.com/Autodeskshore/pkg/command"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulExecWithConfigFile(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execConfig := `{"application": "First Application", "pipeline": "First Pipeline"}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/exec.json"), []byte(execConfig), os.ModePerm)

		// Test

		execCmd := cleanup_command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulExecWithFlag(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execConfig := `{"application": "First Application", "pipeline": "First Pipeline"}`

		viper.Set("payload", execConfig)

		// Test
		execCmd := cleanup_command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestFailureExecWithConfigFileMissingParameter(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execError := "required args key 'pipeline' missing"
		execConfig := `{"application": "First Application"}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/exec.json"), []byte(execConfig), os.ModePerm)

		// Test
		execCmd := cleanup_command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}

func TestFailureExecWithFlagMissingParameter(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execError := "required args key 'pipeline' missing"
		execConfig := `{"application": "First Application"}`

		viper.Set("payload", execConfig)
		command.GetConfigFileOrFlag(deps, "exec", "payload")

		// Test
		execCmd := cleanup_command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}

func TestFailureExecWithConfigFileBadPayload(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execError := "ReadMapCB: expect { or n, but found \x00, error found in #0 byte of ...||..., bigger context ...||..."
		execConfig := `{"application": "First Application",}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/exec.json"), []byte(execConfig), os.ModePerm)

		// Test
		execCmd := cleanup_command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}

func TestFailureExecWithFlagBadPayload(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execError := "required args key 'pipeline' missing\nrequired args key 'application' missing"
		execConfig := `{"application": "First Application",  "pipeline": "First Pipeline" ,,,,,,}`

		viper.Set("payload", execConfig)
		command.GetConfigFileOrFlag(deps, "exec", "payload")

		// Test
		execCmd := cleanup_command.NewExecCommand(deps)
		err := execCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}
