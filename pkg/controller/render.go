package controller

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
	// TODO: For future devx, aggregate errors and return them together.
	projectPath, err := d.Project.GetProjectPath()

	if err != nil {
		return "", err
	}

	renderArgs, err := d.Project.GetRenderArgs()

	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	pipelineJSON, err := d.Renderer.Render(projectPath, renderArgs)

	if err != nil {
		return "", err
	}

	return pipelineJSON, nil
}
