package main_test

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	main "github.com/Autodeskshore/cmd/shore"
)

func TestGetProfileName(t *testing.T) {
	// Given
	givenProfileNameFlag := "test-profile-flag"
	givenProfileNameEnv := "test-profile-env"

	dummyFlaggedCmd := &cobra.Command{}
	dummyFlaggedCmd.Flags().StringP("profile", "P", "", "")
	dummyFlaggedCmd.Flags().Set("profile", givenProfileNameFlag)

	os.Setenv("SHORE_PROFILE", givenProfileNameEnv)
	dummyCmd := &cobra.Command{}

	// Test
	resultingEnvProfileName := main.GetProfileName(dummyCmd)
	os.Unsetenv("SHORE_PROFILE")
	resultingFlagProfileName := main.GetProfileName(dummyFlaggedCmd)
	resultingDefaultProfileName := main.GetProfileName(dummyCmd)

	// Assert
	assert.Equal(t, "default", resultingDefaultProfileName)
	assert.Equal(t, givenProfileNameFlag, resultingFlagProfileName)
	assert.Equal(t, givenProfileNameEnv, resultingEnvProfileName)
}

func TestGetExecutorConfigName(t *testing.T) {
	// Given
	givenExecConfigNameFlag := "test-exec-flag"
	givenExecConfigNameEnv := "test-exec-env"

	dummyFlaggedCmd := &cobra.Command{}
	dummyFlaggedCmd.Flags().StringP("executor-config", "X", "", "")
	dummyFlaggedCmd.Flags().Set("executor-config", givenExecConfigNameFlag)

	os.Setenv("SHORE_EXECUTOR_CONFIG", givenExecConfigNameEnv)
	dummyCmd := &cobra.Command{}

	// Test
	resultingEnvExecConfigName := main.GetExecutorConfigName(dummyCmd)
	resultingFlagExecConfigName := main.GetExecutorConfigName(dummyFlaggedCmd)
	os.Unsetenv("SHORE_EXECUTOR_CONFIG")
	resultingDefaultExecConfigName := main.GetExecutorConfigName(dummyCmd)

	// Assert
	assert.Equal(t, "default", resultingDefaultExecConfigName)
	assert.Equal(t, givenExecConfigNameFlag, resultingFlagExecConfigName)
	assert.Equal(t, givenExecConfigNameEnv, resultingEnvExecConfigName)
}
