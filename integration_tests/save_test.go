package integration_tests

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/Autodeskshore/pkg/command"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulSaveWithConfigFile(t *testing.T) {
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
		saveCmd := command.NewSaveCommand(deps)
		saveCmd.SilenceErrors = true
		saveCmd.SilenceUsage = true
		err := saveCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestFailedSaveMissingParam(t *testing.T) {
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
		saveCmd := command.NewSaveCommand(deps)
		saveCmd.SilenceErrors = true
		saveCmd.SilenceUsage = true
		err := saveCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})
}

func TestSuccessSaveCommandLineRenderParams(t *testing.T) {
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
		saveCmd := command.NewSaveCommand(deps)
		saveCmd.SilenceErrors = true
		saveCmd.SilenceUsage = true
		saveCmd.Flags().Set("render-values", renderConfig)
		err := saveCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}
