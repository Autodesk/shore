package cleanup_test

import (
	"fmt"
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

var testPath string = "/test"

func TestSuccessfulRenderWithConfigFile(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"a": "c", "b": "d"}`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/render.json"), []byte(renderConfig), os.ModePerm)

		pipeline := `
		function(params={})(
			{
				first: params.a,
				second: params.b
			}
		)
		`

		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/cleanup.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := cleanup_command.NewRenderCommand(deps)
		err := renderCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessfulRenderWithoutParams(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		pipeline := `
		function(params={})(
			{
				first: "c",
				second: "d"
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/cleanup.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := cleanup_command.NewRenderCommand(deps)
		err := renderCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestSuccessRenderCommandLineRenderParams(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		renderConfig := `{"a": "c"}`
		viper.Set("values", renderConfig)
		pipeline := `
		function(params={})(
			{
				first: params.a,
				second: "d"
			}
		)
		`
		command.GetConfigFileOrFlag(deps, "render", "values")
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/cleanup.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := cleanup_command.NewRenderCommand(deps)
		err := renderCmd.Execute()

		// Assert
		assert.Nil(t, err)
	})
}

func TestFailedRenderMissingParam(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		expectedErrorMessage := fmt.Sprintf("RUNTIME ERROR: Field does not exist: b\n\t%v/main.pipeline.jsonnet:5:11-26\tobject <anonymous>\n\tDuring manifestation\t\n", testPath)

		renderConfig := `{"a": "c"}`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/render.json"), []byte(renderConfig), os.ModePerm)

		pipeline := `
		function(params={})(
			{
				first: params.a,
				second: params.b
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/cleanup.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := cleanup_command.NewRenderCommand(deps)
		err := renderCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.NotEqual(t, expectedErrorMessage, err.Error())
	})
}

func TestFailedRenderMissingAllRenderFileWithRequiredParams(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
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
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/cleanup.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := cleanup_command.NewRenderCommand(deps)
		err := renderCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, expectedErrorMessage, err.Error())

	})
}

func TestFailedMalformedPipeline(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		expectedErrorMessage := fmt.Sprintf("%s/cleanup/cleanup.pipeline.jsonnet:6:3-4 Unexpected: \")\" while parsing field definition\n\n\t\t)\n\n", testPath)

		pipeline := `
		function(params={})(
			{
				first: params.a,

		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/cleanup.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := cleanup_command.NewRenderCommand(deps)
		err := renderCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, expectedErrorMessage, err.Error())

	})
}

func TestFailedMalformedRenderFile(t *testing.T) {
	integration_tests.SetupTest(t, func(t *testing.T, deps *command.Dependencies) {
		// Given
		expectedErrorMessage := "While parsing config: unexpected end of JSON input"

		renderConfig := `{"a":`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/render.json"), []byte(renderConfig), os.ModePerm)

		pipeline := `
		function(params={})(
			{
				first: params.a,
				second: "b"
			}
		)
		`
		afero.WriteFile(deps.Project.FS, path.Join(testPath, "cleanup/cleanup.pipeline.jsonnet"), []byte(pipeline), os.ModePerm)

		// Test
		renderCmd := cleanup_command.NewRenderCommand(deps)
		err := renderCmd.Execute()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, expectedErrorMessage, err.Error())

	})
}
