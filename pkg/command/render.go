package command

import (
	"fmt"

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

			settingsBytes, err := GetConfigFileOrFlag(d, "render", "values")

			if _, ok := err.(viper.ConfigFileNotFoundError); err != nil && !ok {
				d.Logger.Error("Renderer values could not be loaded, returned an error ", err)
				return err
			}

			pipeline, err := Render(d, settingsBytes)

			if err != nil {
				return err
			}

			fmt.Println(pipeline)
			return nil
		},
	}

	cmd.Flags().StringP("values", "r", "", "A JSON string for the render. If not provided the render.[json/yml/yaml] file is used.")
	viper.BindPFlag("values", cmd.Flags().Lookup("values"))

	return cmd
}

// Render - Using a Project & Renderer, renders the pipeline.
func Render(d *Dependencies, settings []byte) (string, error) {
	// TODO: For future DevX, aggregate errors and return them together.
	d.Logger.Info("Render function started")

	d.Logger.Debug("GetProjectPath")
	projectPath, err := d.Project.GetProjectPath()

	if err != nil {
		d.Logger.Error("GetProjectPath returned an error ", err)
		return "", err
	}

	d.Logger.Debug("GetProjectPath returned ", projectPath)

	// A bit of a hack, rather change this to an object later on.
	renderArgs := string(settings)

	d.Logger.Debug("Args returned:\n", renderArgs)

	d.Logger.Info("calling Renderer.Render with projectPath ", projectPath, " and renderArgs ", renderArgs)
	pipelineJSON, err := d.Renderer.Render(projectPath, renderArgs)

	if err != nil {
		d.Logger.Error("Renderer.Render returned an error ", err)
		return "", err
	}

	d.Logger.Debug("Renderer.Render returned successfully")

	return pipelineJSON, nil
}
