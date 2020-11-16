package controller

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewSaveCommand - Using a Project, Renderer & Backend, renders and saves a pipeline.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewSaveCommand(d *Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "save",
		Short: "save the pipeline",
		Long:  "Using the main file configured by the renderer save the pipeline (or pipelines)",
		RunE: func(cmd *cobra.Command, args []string) error {
			d.Logger.Info("Calling save pipeline")
			pipeline, err := Render(d)

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
}
