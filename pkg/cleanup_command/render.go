package cleanup_command

import (
	"errors"
	"fmt"

	"github.com/Autodeskshore/pkg/command"
	"github.com/Autodeskshore/pkg/config"
	"github.com/Autodeskshore/pkg/renderer"
	"github.com/spf13/cobra"
)

// NewRenderCommand - A Cobra wrapper for the common Render function.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewRenderCommand(d *command.Dependencies) *cobra.Command {
	var values string

	cmd := &cobra.Command{
		Use:   "render",
		Short: "render the cleanup pipeline",
		Long: `Renders the "cleanup.pipeline.jsonnet" file.
This helper utility command is used to debug issues when the "cleanup" pipeline doesn't render correctly.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			settingsBytes, err := config.LoadConfig(d.Project, values, "cleanup/render")

			var confErr *config.FileConfErr

			if err != nil && !errors.As(errors.Unwrap(err), &confErr) {
				return err
			}

			pipeline, err := command.Render(d, settingsBytes, renderer.CleanUpFileName)
			if err != nil {
				return err
			}

			fmt.Println(pipeline)
			return nil
		},
	}

	cmd.Flags().StringVarP(&values, "values", "r", "", "A JSON string for the render. If not provided the render.[json/yml/yaml] file is used.")

	return cmd
}
