package command

import (
	"fmt"

	// Feels a bit weird, maybe move the TestsConfig object out?

	"github.com/Autodeskshore/pkg/shore_testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
)

/*
Controller E2E is a test suite implementation to allow developers to test their pipelines in the intended backend.
Testing implementation should be defined on the `backend` level.
*/

// NewTestRemoteCommand - Using a Project, Renderer & Backend runs a test suite defined in a config file and outputs the results to the customer.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang/etc...) and backends (Spinnaker/Tekton/ArgoCD/JenkinsX/etc...)
func NewTestRemoteCommand(d *Dependencies) *cobra.Command {
	var testNames []string
	var stringifyNonScalars bool

	cmd := &cobra.Command{
		Use:   "test-remote",
		Short: "Run the test suite on a remotely saved pipeline",
		Long:  "Using the E2E.yaml file run a full test-suite on the pipeline stored in a specific backend",
		RunE: func(cmd *cobra.Command, args []string) error {
			testSettingsBytes, err := GetConfigFileOrFlag(d, "E2E", "")

			if err != nil {
				return err
			}

			var testConfig shore_testing.TestsConfig
			if err := jsoniter.Unmarshal(testSettingsBytes, &testConfig); err != nil {
				return err
			}

			if len(testNames) > 0 {
				if err := verifyTestExist(testNames, testConfig); err != nil {
					return err
				}
				testConfig.Ordering = testNames
			} else if testConfig.Ordering != nil && len(testConfig.Ordering) > 0 {
				if err := verifyTestExist(testConfig.Ordering, testConfig); err != nil {
					return err
				}
			}

			d.Logger.Debug("Stringify is ", stringifyNonScalars)

			err = d.Backend.TestPipeline(testConfig, func() {}, stringifyNonScalars)

			if err != nil {
				return err
			}

			fmt.Println("Test Passed!")

			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&testNames, "test-names", "t", []string{}, "An array of tests that will be ran. Preserves order.")
	cmd.Flags().BoolVarP(&stringifyNonScalars, "stringify", "y", true, "Stringifies the non scalar parameters to SpinCli")

	return cmd
}

func verifyTestExist(testNames []string, testConfig shore_testing.TestsConfig) error {
	for _, testName := range testNames {
		if _, ok := testConfig.Tests[testName]; !ok {
			return fmt.Errorf("the provided E2E configuration does not contain [%s] as a test", testName)
		}
	}
	return nil
}
