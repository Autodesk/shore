package command

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Autodesk/shore/pkg/config"
	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
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

			s := spinner.New(spinner.CharSets[9], 1000*time.Millisecond)
			s.Writer = color.Error
			s.Suffix = " Deleting spinnaker pipelines, this may take a few moments (depending on Internet traffic!)\n"
			s.Start()
			res, err := d.Backend.DeletePipeline(pipeline, dryRun)
			s.Stop()

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
