package cleanup_command

import (
	"fmt"

	"github.com/Autodeskshore/pkg/command"
	"github.com/Autodeskshore/pkg/renderer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRenderCommand - A Cobra wrapper for the common Render function.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewRenderCommand(d *command.Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render",
		Short: "render the cleanup pipeline",
		Long: `Renders the "cleanup.pipeline.jsonnet" file.
This helper utility command is used to debug issues when the "cleanup" pipeline doesn't render correctly.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			settingsBytes, err := command.GetConfigFileOrFlag(d, "cleanup/render", "values")

			if _, ok := err.(viper.ConfigFileNotFoundError); err != nil && !ok {
				d.Logger.Error("Renderer values could not be loaded, returned an error ", err)
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

	cmd.Flags().StringP("values", "r", "", "A JSON string for the render. If not provided the render.[json/yml/yaml] file is used.")
	viper.BindPFlag("values", cmd.Flags().Lookup("values"))

	return cmd
}
