package cleanup_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/Autodesk/shore/integration_tests"
	"github.com/Autodesk/shore/pkg/cleanup_command"
	"github.com/Autodesk/shore/pkg/command"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulSaveWithConfigFile(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"application": "First Application", "pipeline": "First Pipeline"}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/render.json"), []byte(renderConfig), os.ModePerm)

		pipeline := `
		function(params={})(
			{
				application: params.application,
				name: params.pipeline
			}
		)
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/cleanup.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test

		saveCmd := cleanup_command.NewSaveCommand(deps)
		saveCmd.SilenceErrors = true
		saveCmd.SilenceUsage = true
		err := saveCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestFailedSaveMissingParam(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		expectedErrorMessage := fmt.Sprintf("RUNTIME ERROR: Field does not exist: pipeline\n\t%s/cleanup/cleanup.pipeline.jsonnet:5:11-26\tobject <anonymous>\n\tField \"name\"\t\n\tDuring manifestation\t\n", testPath)

		renderConfig := `{"application":"Fourth Application"}`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/render.json"), []byte(renderConfig), os.ModePerm)

		pipeline := `
		function(params={})(
			{
				application: params.application,
				name: params.pipeline
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/cleanup.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		saveCmd := cleanup_command.NewSaveCommand(deps)
		saveCmd.SilenceErrors = true
		saveCmd.SilenceUsage = true
		err := saveCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})
}

func TestSuccessSaveCommandLineRenderParams(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
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
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/cleanup.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		saveCmd := cleanup_command.NewSaveCommand(deps)
		saveCmd.SilenceErrors = true
		saveCmd.SilenceUsage = true

		saveCmd.Flags().Set("render-values", renderConfig)
		err := saveCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}
