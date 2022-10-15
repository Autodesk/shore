package cleanup_test

import (
	"os"
	"path"
	"testing"

	"github.com/Autodesk/shore/integration_tests"
	"github.com/Autodesk/shore/pkg/cleanup_command"
	"github.com/Autodesk/shore/pkg/command"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulExecWithConfigFile(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execConfig := `{"application": "First Application", "pipeline": "First Pipeline"}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/exec.json"), []byte(execConfig), os.ModePerm)

		// Test

		execCmd := cleanup_command.NewExecCommand(deps)
		execCmd.SilenceErrors = true
		execCmd.SilenceUsage = true
		err := execCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulExecWithFlag(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execConfig := `{"application": "First Application", "pipeline": "First Pipeline"}`

		// Test
		execCmd := cleanup_command.NewExecCommand(deps)
		execCmd.SilenceErrors = true
		execCmd.SilenceUsage = true
		execCmd.Flags().Set("payload", execConfig)
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
		execCmd.SilenceErrors = true
		execCmd.SilenceUsage = true
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

		// Test
		execCmd := cleanup_command.NewExecCommand(deps)
		execCmd.SilenceErrors = true
		execCmd.SilenceUsage = true
		execCmd.Flags().Set("payload", execConfig)
		err := execCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, execError, err.Error())
	})
}

func TestFailureExecWithConfigFileBadPayload(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execConfig := `{"application": "First Application",}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/exec.json"), []byte(execConfig), os.ModePerm)

		// Test
		execCmd := cleanup_command.NewExecCommand(deps)
		execCmd.SilenceErrors = true
		execCmd.SilenceUsage = true
		err := execCmd.Execute()

		// Assert
		assert.Error(t, err)
	})
}

func TestFailureExecWithFlagBadPayload(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		execConfig := `{"application": "First Application",  "pipeline": "First Pipeline" ,,,,,,}`

		// Test
		execCmd := cleanup_command.NewExecCommand(deps)
		execCmd.SilenceErrors = true
		execCmd.SilenceUsage = true
		execCmd.Flags().Set("payload", execConfig)
		err := execCmd.Execute()

		// Assert
		assert.Error(t, err)
	})
}
