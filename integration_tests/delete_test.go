package integration_tests

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/Autodesk/shore/pkg/command"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulDeleteWithConfigFile(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"application": "First Application", "pipeline": "First Pipeline"}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "render.json"), []byte(renderConfig), os.ModePerm)

		pipeline := `
		function(params={})(
			{
				application: params.application,
				name: params.pipeline
			}
		)
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "main.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		deleteCmd := command.NewDeleteCommand(deps)
		deleteCmd.SilenceErrors = true
		deleteCmd.SilenceUsage = true
		err := deleteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestFailedDeleteMissingParam(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		expectedErrorMessage := fmt.Sprintf("RUNTIME ERROR: Field does not exist: pipeline\n\t%v/main.pipeline.jsonnet:5:11-26\tobject <anonymous>\n\tField \"name\"\t\n\tDuring manifestation\t\n", testPath)

		renderConfig := `{"application":"Fourth Application"}`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "render.json"), []byte(renderConfig), os.ModePerm)

		pipeline := `
		function(params={})(
			{
				application: params.application,
				name: params.pipeline
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "main.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		deleteCmd := command.NewDeleteCommand(deps)
		deleteCmd.SilenceErrors = true
		deleteCmd.SilenceUsage = true
		err := deleteCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})
}

func TestSuccessDeleteCommandLineDryRun(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"application":"Third Application"}`
		pipeline := `
		function(params={})(
			{
				application: params.application,
				name: "Third Pipeline"
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "main.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		deleteCmd := command.NewDeleteCommand(deps)
		deleteCmd.SilenceErrors = true
		deleteCmd.SilenceUsage = true
		deleteCmd.Flags().Set("render-values", renderConfig)
		deleteCmd.Flags().Set("dry-run", "true")
		err := deleteCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}
