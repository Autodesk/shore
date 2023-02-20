package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/Autodesk/shore/pkg/config"
	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/spf13/cobra"
)

// NewDeleteCommand - Using a Project, Renderer & Backend, renders and saves a pipeline.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewDeleteCommand(d *Dependencies) *cobra.Command {
	var renderVals string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete the pipeline",
		Long:  "Using the main file configured by the renderer delete the pipeline (or pipelines)",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsBytes, err := config.LoadConfig(d.Project, renderVals, "render")

			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}

			pipeline, err := Render(d, settingsBytes, renderer.MainFileName)

			if err != nil {
				return err
			}

			d.Logger.Info("Calling Backend.DeletePipeline")
			res, err := d.Backend.DeletePipeline(pipeline, dryRun)

			if err != nil {
				d.Logger.Warnf("Delete pipeline returned an error: %v", err)
				return err
			}

			d.Logger.Info("Backend.DeletePipeline returned")
			fmt.Println(res)
			return nil
		},
	}

	cmd.Flags().StringVarP(&renderVals, "render-values", "r", "", "A JSON string for the render. If not provided the render.[json/yml/yaml] file is used.")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "list pipelines to be deleted - dry run")

	return cmd
}
