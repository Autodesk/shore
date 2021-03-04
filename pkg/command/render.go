package command

import (
	"fmt"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRenderCommand - A Cobra wrapper for the common Render function.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewRenderCommand(d *Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render",
		Short: "render a pipeline",
		Long:  "Walk through the `pipelines` directory, renderer the pipelines and output to STDOUT",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigName("render")

			var valuesErr error

			if viper.IsSet("values") {
				valuesErr = viper.ReadConfig(strings.NewReader(viper.GetString("values")))
			} else {
				valuesErr = viper.ReadInConfig()
			}

			if valuesErr != nil {
				d.Logger.Error("Failed to load values.")
				return valuesErr
			}

			pipeline, err := Render(d)

			if err != nil {
				return err
			}

			fmt.Println(pipeline)
			return nil
		},
	}

	cmd.Flags().StringP("values", "r", "", "A JSON string for the render. If not provided the render.[json/jyml/yaml] file is used.")
	viper.BindPFlag("values", cmd.Flags().Lookup("values"))

	return cmd
}

// Render - Using a Project & Renderer, renders the pipeline.
func Render(d *Dependencies) (string, error) {
	// TODO: For future DevX, aggregate errors and return them together.
	d.Logger.Info("Render function started")

	d.Logger.Debug("GetProjectPath")
	projectPath, err := d.Project.GetProjectPath()

	if err != nil {
		d.Logger.Error("GetProjectPath returned an error ", err)
		return "", err
	}

	d.Logger.Debug("GetProjectPath returned ", projectPath)

	renderArgs := ""

	values := viper.AllSettings()
	valuesBytes, err := jsoniter.Marshal(values)

	if err != nil && !os.IsNotExist(err) {
		d.Logger.Error("Renderer values could not be loaded, returned an error ", err)
		return "", err
	}

	// A bit of a hack, rather change this to an object later on.
	renderArgs = string(valuesBytes)

	d.Logger.Debug("Args returned:\n", renderArgs)

	d.Logger.Info("calling Renderer.Render with projectPath ", projectPath, "and renderArgs ", renderArgs)
	pipelineJSON, err := d.Renderer.Render(projectPath, renderArgs)

	if err != nil {
		d.Logger.Error("Renderer.Render returned an error ", err)
		return "", err
	}

	d.Logger.Debug("Renderer.Render returned successfully")

	return pipelineJSON, nil
}
