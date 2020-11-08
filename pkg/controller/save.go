package controller

import (
	"log"

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
			pipeline, err := Render(d)

			if err != nil {
				return err
			}

			res, err := d.Backend.SavePipeline(pipeline)

			if err != nil {
				return err
			}

			log.Println(res)
			return nil
		},
	}
}
