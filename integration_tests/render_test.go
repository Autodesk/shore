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

func TestSuccessfulRenderWithConfigFile(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"a": "c", "b": "d"}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "render.json"), []byte(renderConfig), os.ModePerm)

		pipeline := `
		function(params={})(
			{
				first: params.a,
				second: params.b
			}
		)
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "main.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := command.NewRenderCommand(deps)
		renderCmd.SilenceErrors = true
		renderCmd.SilenceUsage = true
		err := renderCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulRenderWithoutParams(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		pipeline := `
		function(params={})(
			{
				first: "c",
				second: "d"
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "main.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := command.NewRenderCommand(deps)
		renderCmd.SilenceErrors = true
		renderCmd.SilenceUsage = true
		err := renderCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessRenderCommandLineRenderParams(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"a": "c"}`
		pipeline := `
		function(params={})(
			{
				first: params.a,
				second: "d"
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "main.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := command.NewRenderCommand(deps)
		renderCmd.SilenceErrors = true
		renderCmd.SilenceUsage = true
		renderCmd.Flags().Set("values", renderConfig)
		err := renderCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestFailedRenderMissingParam(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		expectedErrorMessage := fmt.Sprintf("RUNTIME ERROR: Field does not exist: b\n\t%v/main.pipeline.jsonnet:5:11-26\tobject <anonymous>\n\tDuring manifestation\t\n", testPath)

		renderConfig := `{"a": "c"}`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "render.json"), []byte(renderConfig), os.ModePerm)

		pipeline := `
		function(params={})(
			{
				first: params.a,
				second: params.b
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "main.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := command.NewRenderCommand(deps)
		renderCmd.SilenceErrors = true
		renderCmd.SilenceUsage = true
		err := renderCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.NotEqual(t, expectedErrorMessage, err.Error())
	})
}

func TestFailedRenderMissingAllRenderFileWithRequiredParams(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		expectedErrorMessage := "<top-level-arg:params>:1:1 Unexpected end of file\n\n\n\n"

		pipeline := `
		function(params={})(
			{
				a: params.a,
				name: params.pipeline
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "main.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := command.NewRenderCommand(deps)
		renderCmd.SilenceErrors = true
		renderCmd.SilenceUsage = true
		err := renderCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, expectedErrorMessage, err.Error())

	})
}

func TestFailedMalformedPipeline(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		expectedErrorMessage := fmt.Sprintf("%v/main.pipeline.jsonnet:6:3-4 Unexpected: \")\" while parsing field definition\n\n\t\t)\n\n", testPath)

		pipeline := `
		function(params={})(
			{
				first: params.a,

		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "main.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := command.NewRenderCommand(deps)
		renderCmd.SilenceErrors = true
		renderCmd.SilenceUsage = true
		err := renderCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, expectedErrorMessage, err.Error())

	})
}

func TestFailedMalformedRenderFile(t *testing.T) {
	SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"a":`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "render.json"), []byte(renderConfig), os.ModePerm)

		pipeline := `
		function(params={})(
			{
				first: params.a,
				second: "b"
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "main.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := command.NewRenderCommand(deps)
		renderCmd.SilenceErrors = true
		renderCmd.SilenceUsage = true
		err := renderCmd.Execute()

		// Assert
		assert.Error(t, err)
	})
}
