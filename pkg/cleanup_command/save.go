package cleanup_command

import (
	"fmt"

	"github.com/Autodesk/shore/pkg/command"
	"github.com/Autodesk/shore/pkg/config"
	"github.com/Autodesk/shore/pkg/renderer"
	"github.com/spf13/cobra"
)

// NewSaveCommand - Using a Project, Renderer & Backend, renders and saves a pipeline.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewSaveCommand(d *command.Dependencies) *cobra.Command {
	var renderValues string

	cmd := &cobra.Command{
		Use:   "save",
		Short: "save the cleanup pipeline",
		Long: `Save the cleanup pipeline to the selected backend.
Help in developing and debugging cleanup pipelines in a live environment.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsBytes, _ := config.LoadConfig(d.Project, renderValues, "cleanup/render")

			pipeline, err := command.Render(d, settingsBytes, renderer.CleanUpFileName)

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

	cmd.Flags().StringVarP(&renderValues, "render-values", "r", "", "A JSON string for the render. If not provided the render.[json/yml/yaml] file is used.")

	return cmd
}
