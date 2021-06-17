package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

/*
Controller E2E is a test suite implementation to allow developers to test their pipelines in the intended backend.
Testing implementation should be defined on the `backend` level.
*/

// NewTestRemoteCommand - Using a Project, Renderer & Backend runs a test suite defined in a config file and outputs the results to the customer.
// Abstraction for different configuration languages (I.E. Jsonnet/HCL/CUELang/etc...) and backends (Spinnaker/Tekton/ArgoCD/JenkinsX/etc...)
func NewTestRemoteCommand(d *Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "test-remote",
		Short: "Run the test suite on a remotely saved pipeline",
		Long:  "Using the E2E.yaml file run a full test-suite on the pipeline stored in a specific backend",
		RunE: func(cmd *cobra.Command, args []string) error {
			testSettingsBytes, err := GetConfigFileOrFlag(d, "E2E", "")

			if err != nil {
				return err
			}

			// A bit of a hack, rather change this to an object later on.
			testConfig := string(testSettingsBytes)

			err = d.Backend.TestPipeline(testConfig, func() {})

			if err != nil {
				return err
			}

			fmt.Println("Test Passed!")

			return nil
		},
	}
}
