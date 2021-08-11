package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/Autodeskshore/pkg/renderer"
	"github.com/spf13/cobra"
)

// NewSaveCommand - Using a Project, Renderer & Backend, renders and saves a pipeline.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewSaveCommand(d *Dependencies) *cobra.Command {
	var renderVals string

	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save the pipeline",
		Long:  "Using the main file configured by the renderer save the pipeline (or pipelines)",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsBytes, err := GetConfigFileOrFlag(d, "render", renderVals)

			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}

			pipeline, err := Render(d, settingsBytes, renderer.MainFileName)

			if err != nil {
				return err
			}

			d.Logger.Info("Calling Backend.SavePipeline")
			res, err := d.Backend.SavePipeline(pipeline)

			if err != nil {
				d.Logger.Warn("Save pipeline returned an error", err)
				return err
			}

			d.Logger.Info("Backend.SavePipeline returned")
			fmt.Println(res)
			return nil
		},
	}

	cmd.Flags().StringVarP(&renderVals, "render-values", "r", "", "A JSON string for the render. If not provided the render.[json/yml/yaml] file is used.")

	return cmd
}
