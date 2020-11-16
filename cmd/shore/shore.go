package main

import (
	"fmt"
	"os"

	"github.com/Autodesk/shore/pkg/backend/spinnaker"
	"github.com/Autodesk/shore/pkg/controller"
	"github.com/Autodesk/shore/pkg/project"
	"github.com/Autodesk/shore/pkg/renderer/jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


var logVerbosity int
var logger *logrus.Logger

var rootCmd = &cobra.Command{
	Use:           "shore",
	Short:         "Shore - Pipeline Framework",
	Long:          "A Pipeline development framework for integrated pipelines.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(*cobra.Command, []string) {
		logLevel := logrus.WarnLevel + logrus.Level(logVerbosity)
		logger.SetLevel(logLevel)
		logger.SetFormatter(&logrus.TextFormatter{})
	},
}

func init() {
	// TODO: Add global validations to init.
	// cobra.OnInitialize()
	viper.AutomaticEnv()
	fs := afero.NewOsFs()
	logger = logrus.New()

	commonDependencies := &controller.Dependencies{
		Project:  project.NewShoreProject(fs, logger),
		Renderer: jsonnet.NewRenderer(fs, logger),
		Backend:  spinnaker.NewClient(logger),
		Logger:   logger,
	}

	rootCmd.PersistentFlags().CountVarP(&logVerbosity, "verbose", "v", "Logging verbosity")

	rootCmd.AddCommand(controller.NewRenderCommand(commonDependencies))
	rootCmd.AddCommand(controller.NewSaveCommand(commonDependencies))
	rootCmd.AddCommand(controller.NewExecCommand(commonDependencies))
	rootCmd.AddCommand(controller.NewTestRemoteCommand(commonDependencies))
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
