package main

import (
	"fmt"
	"os"

	"github.com/Autodesk/shore/pkg/cleanup_command"
	"github.com/Autodesk/shore/pkg/command"
	"github.com/Autodesk/shore/pkg/config"
	"github.com/Autodesk/shore/pkg/project"
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

var commonDependencies *command.Dependencies

var rootCmd = &cobra.Command{
	Use:           "shore",
	Short:         "Shore - Pipeline Framework",
	Long:          "A Pipeline development framework for integrated pipelines.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logLevel := logrus.WarnLevel + logrus.Level(logVerbosity)
		logger.SetLevel(logLevel)
		logger.SetFormatter(&logrus.TextFormatter{})

		if cmd.Name() == "help" {
			return nil // No need to do anything, just printing help
		}

		if err := commonDependencies.Load(); err != nil {
			return err
		}

		commonDependencies.ShoreConfigOpts.ProfileName = GetProfileName(cmd)
		commonDependencies.ShoreConfigOpts.ExecutorConfigName = GetExecutorConfigName(cmd)

		logger.Debug("Profile set to - ", commonDependencies.ShoreConfigOpts.ProfileName)
		logger.Debug("Executor configuration set to - ", commonDependencies.ShoreConfigOpts.ExecutorConfigName)
		return nil
	},
}

// GetProfileName - Gets the Profile name based on the env var or flag.
func GetProfileName(cmd *cobra.Command) string {
	return getConfigName(cmd, "profile", "SHORE_PROFILE")
}

// GetExecutorConfigName - Gets the Backend config name based on the env var or flag.
func GetExecutorConfigName(cmd *cobra.Command) string {
	return getConfigName(cmd, "executor-config", "SHORE_EXECUTOR_CONFIG")
}

func getConfigName(cmd *cobra.Command, flagName string, envVar string) string {
	configName := "default"
	flagValue, err := cmd.Flags().GetString(flagName)
	envValue := os.Getenv(envVar) // Either a non-empty string, or an empty string

	if len(envValue) > 0 {
		configName = envValue
	}
	if len(flagValue) > 0 && err == nil {
		configName = flagValue
	}

	return configName
}

func init() {
	// TODO: Add global validations to init.
	// cobra.OnInitialize()
	fs := afero.NewOsFs()
	logger = logrus.New()

	commonDependencies = &command.Dependencies{
		Project: project.NewShoreProject(fs, logger),
		Logger:  logger,
		ShoreConfigOpts: config.ShoreConfigOpts{
			ProfileName:        "default",
			ExecutorConfigName: "default",
		},
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
