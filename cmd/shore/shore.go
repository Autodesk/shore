package main

import (
	"fmt"
	"os"

	"github.com/Autodesk/shore/pkg/backend/spinnaker"
	"github.com/Autodesk/shore/pkg/controller"
	"github.com/Autodesk/shore/pkg/project"
	"github.com/Autodesk/shore/pkg/renderer/jsonnet"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:          "shore",
	Short:        "Shore - Pipeline Framework",
	Long:         "A Pipeline development framework for integrated pipelines.",
	SilenceUsage: true,
}

func init() {
	// TODO: Add global validations to init.
	// cobra.OnInitialize()
	viper.AutomaticEnv()
	fs := afero.NewOsFs()

	commonDependencies := &controller.Dependencies{
		Project:  project.NewShoreProject(fs),
		Renderer: jsonnet.NewRenderer(fs),
		Backend:  spinnaker.NewClient(),
	}

	rootCmd.AddCommand(controller.NewRenderCommand(commonDependencies))
	rootCmd.AddCommand(controller.NewSaveCommand(commonDependencies))
	rootCmd.AddCommand(controller.NewExecCommand(commonDependencies))
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	execute()
}
