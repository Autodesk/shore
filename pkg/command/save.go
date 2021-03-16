package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewSaveCommand - Using a Project, Renderer & Backend, renders and saves a pipeline.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang)
func NewSaveCommand(d *Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save",
		Short: "save the pipeline",
		Long:  "Using the main file configured by the renderer save the pipeline (or pipelines)",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigName("render")

			d.Logger.Info("Calling save pipeline")

			var valuesErr error

			if viper.IsSet("values") {
				valuesErr = viper.ReadConfig(strings.NewReader(viper.GetString("render-values")))
			} else {
				valuesErr = viper.ReadInConfig()
			}

			if valuesErr != nil {
				if _, ok := valuesErr.(viper.ConfigFileNotFoundError); ok {
					d.Logger.Warn(valuesErr)
				} else {
					d.Logger.Error("Failed to load values")
					return valuesErr
				}

			}

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

	cmd.Flags().StringP("render-values", "r", "", "A JSON string for the render. If not provided the render.[json/jyml/yaml] file is used.")
	viper.BindPFlag("render-values", cmd.Flags().Lookup("render-values"))

	return cmd
}
