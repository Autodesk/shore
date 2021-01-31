package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRenderCommand - A Cobra wrapper for the common Render function.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewRenderCommand(d *Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "render",
		Short: "render a pipeline",
		Long:  "Walk through the `pipelines` directory, renderer the pipelines and output to STDOUT",
		RunE: func(cmd *cobra.Command, args []string) error {
			pipeline, err := Render(d)

			if err != nil {
				return err
			}

			fmt.Println(pipeline)
			return nil
		},
	}
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

	d.Logger.Debug("GetRenderArgs")
	renderArgs, err := d.Project.GetRenderArgs()

	if err != nil && !os.IsNotExist(err) {
		d.Logger.Error("GetRenderArgs returned an error ", err)
		return "", err
	}

	d.Logger.Debug("GetRenderArgs returned ", renderArgs)

	d.Logger.Info("calling Renderer.Render with projectPath ", projectPath, "and renderArgs ", renderArgs)
	pipelineJSON, err := d.Renderer.Render(projectPath, renderArgs)

	if err != nil {
		d.Logger.Error("Renderer.Render returned an error ", err)
		return "", err
	}

	d.Logger.Debug("Renderer.Render returned successfully")

	return pipelineJSON, nil
}
