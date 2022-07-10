package main

import (
	"fmt"
	"os"

	"github.com/Autodeskshore/pkg/backend/spinnaker"
	"github.com/Autodeskshore/pkg/cleanup_command"
	"github.com/Autodeskshore/pkg/command"
	"github.com/Autodeskshore/pkg/project"
	"github.com/Autodeskshore/pkg/renderer/jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// Version - Shore CLI version
// This variable is populated during compile time with a variable substitution.
// The variable should be a `const`, but `ldflags` can only operate on `var+string`.
var version = "local"

var logVerbosity int
var logger *logrus.Logger

var rootCmd = &cobra.Command{
	Use:           "shore",
	Short:         "Shore - Pipeline Framework",
	Long:          "A Pipeline development framework for integrated pipelines.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       version,
	PersistentPreRun: func(*cobra.Command, []string) {
		logLevel := logrus.WarnLevel + logrus.Level(logVerbosity)
		logger.SetLevel(logLevel)
		logger.SetFormatter(&logrus.TextFormatter{})
	},
}

func init() {
	// TODO: Add global validations to init.
	// cobra.OnInitialize()
	fs := afero.NewOsFs()
	logger = logrus.New()

	commonDependencies := &command.Dependencies{
		Project:  project.NewShoreProject(fs, logger),
		Renderer: jsonnet.NewRenderer(fs, logger),
		Backend:  spinnaker.NewClient(logger),
		Logger:   logger,
	}

	rootCmd.PersistentFlags().CountVarP(&logVerbosity, "verbose", "v", "Logging verbosity")

	rootCmd.AddCommand(command.NewProjectCommand(commonDependencies))
	rootCmd.AddCommand(command.NewRenderCommand(commonDependencies))
	rootCmd.AddCommand(command.NewDiffCommand(commonDependencies))
	rootCmd.AddCommand(command.NewSaveCommand(commonDependencies))
	rootCmd.AddCommand(command.NewExecCommand(commonDependencies, "exec"))
	rootCmd.AddCommand(command.NewTestRemoteCommand(commonDependencies))
	rootCmd.AddCommand(cleanup_command.NewCleanupCommand(commonDependencies))
	// Make the version easily parsable when invoking `shore --version`
	rootCmd.SetVersionTemplate("{{.Version}}\n")
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
