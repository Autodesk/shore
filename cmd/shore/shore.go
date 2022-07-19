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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logLevel := logrus.WarnLevel + logrus.Level(logVerbosity)
		logger.SetLevel(logLevel)
		logger.SetFormatter(&logrus.TextFormatter{})

		profileName := command.GetProfileName(cmd)
		ExecConfigName := command.GetExecutorConfigName(cmd)

		logger.Debug("Profile set to - ", profileName)
		logger.Debug("Executor configuration set to - ", ExecConfigName)
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
	// "default" should not be set explicitly on the command - it will be set in getConfigName.
	rootCmd.PersistentFlags().StringP("executor-config", "X", os.Getenv("SHORE_EXECUTOR_CONFIG"),
		"The backend configuration name to use. Can also be set by $SHORE_EXECUTOR_CONFIG environment variable. Priority is: env variable, cli args, default.")
	//'p' is used for 'payload' used by exec command. 'l' for load profile?
	rootCmd.PersistentFlags().StringP("profile", "P", os.Getenv("SHORE_PROFILE"),
		"The profile to use. Can also be set by $SHORE_PROFILE environment variable. Priority is: env variable, cli args, default.")

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
