package command

import (
	"errors"
	"fmt"

	"github.com/Autodeskshore/pkg/renderer"
	"github.com/spf13/cobra"
)

// NewRenderCommand - A Cobra wrapper for the common Render function.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewRenderCommand(d *Dependencies) *cobra.Command {
	var renderValues string

	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render the pipeline",
		Long: `Render the "main.pipeline.jsonnet" file.
Automatically reads libraries from "vendor/". The Jsonnet-Bundler default path for libraries`,
		RunE: func(cmd *cobra.Command, args []string) error {

			settingsBytes, err := GetConfigFileOrFlag(d, "render", renderValues)

			var confErr *DefaultConfErr

			if err != nil && !errors.As(errors.Unwrap(err), &confErr) {
				return err
			}

			pipeline, err := Render(d, settingsBytes, renderer.MainFileName)

			if err != nil {
				return err
			}

			fmt.Println(pipeline)
			return nil
		},
	}

	cmd.Flags().StringVarP(&renderValues, "values", "r", "", "A JSON string for the render. If not provided the render.[json/yml/yaml] file is used.")

	return cmd
}

// Render - Using a Project & Renderer, renders the pipeline.
func Render(d *Dependencies, settings []byte, renderType renderer.RenderType) (string, error) {
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
	pipelineJSON, err := d.Renderer.Render(projectPath, renderArgs, renderType)

	if err != nil {
		d.Logger.Error("Renderer.Render returned an error ", err)
		return "", err
	}

	d.Logger.Debug("Renderer.Render returned successfully")

	return pipelineJSON, nil
}
