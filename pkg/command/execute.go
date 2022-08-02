package command

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Autodeskshore/pkg/config"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewExecCommand - Using a Project, Renderer & Backend, executes a pipeline pipeline.
func NewExecCommand(d *Dependencies, configPath string) *cobra.Command {
	var withWait bool
	var withSilent bool
	var waitTimeout int
	var withPayload string
	var stringifyNonScalars bool

	cmd := &cobra.Command{
		Use:   "exec",
		Short: "Executes the pipeline",
		Long:  "Executes the selected pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsBytes, err := config.LoadConfig(d.Project, withPayload, configPath)

			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}

			// A bit of a hack, rather change this to an object later on.
			execArgs := string(settingsBytes)

			d.Logger.Debug("Stringify is ", stringifyNonScalars)

			d.Logger.Debug("Calling `Backend.ExecutePipeline`")
			refID, res, err := d.Backend.ExecutePipeline(execArgs, stringifyNonScalars)

			if err != nil {
				return err
			}

			if !withWait {
				// Skip if the output should be silent.
				if !withSilent {
					fmt.Println(res)
				}

				return nil
			}

			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Writer = color.Error
			s.Suffix = fmt.Sprintf(" Waiting for pipeline to finish executing (%d Seconds)", waitTimeout)
			s.Start() // Start the spinner
			execDetails, _, err := d.Backend.WaitForPipelineToFinish(refID, waitTimeout)
			s.Stop() // Stop the spinner

			if err != nil {
				return err
			}
			// Return early if the output should be silent.
			if withSilent {
				return nil
			}

			fmt.Println(execDetails)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&withWait, "wait", "w", false, "Wait for the pipeline to finish execution")
	cmd.Flags().BoolVarP(&withSilent, "silent", "s", false, "Do not print JSON response to STDOUT")
	cmd.Flags().IntVarP(&waitTimeout, "timeout", "t", 60, "how long to wait (Seconds) for the pipeline to finish in Seconds. Yes Seconds.")
	cmd.Flags().StringVarP(&withPayload, "payload", "p", "", "A JSON payload string. If not provided the exec.[json/yml/yaml] file is used.")
	cmd.Flags().BoolVarP(&stringifyNonScalars, "stringify", "y", true, "Stringifies the non scalar parameters to SpinCli")

	return cmd
}
