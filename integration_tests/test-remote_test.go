package integration_tests

import (
	"os"
	"path"
	"testing"

	"github.com/Autodeskshore/pkg/command"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulRemoteTestWithConfigFile(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		e2eConfig := `
		{
			"application": "cosv3-state-buckets",
			"pipeline": "cosv3-dynamodb-table",
			"tests": {
				"Test Success": {
					"assertions": {
						"testedname": {
							"expected_output": {
								"test": "123"
							},
							"expected_status": "succeeded"
						}
					}
				}
			}
		}
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "E2E.json"), []byte(e2eConfig), os.ModePerm)

		// Test
		testRemoteCmd := command.NewTestRemoteCommand(deps)
		err := testRemoteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}
