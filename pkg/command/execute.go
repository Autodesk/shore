package command

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewExecCommand - Using a Project, Renderer & Backend, executes a pipeline pipeline.
func NewExecCommand(d *Dependencies) *cobra.Command {
	var withWait bool
	var withSilent bool
	var waitTimeout int

	cmd := &cobra.Command{
		Use:   "exec",
		Short: "executes the pipeline",
		Long:  "Executes the selected pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			d.Logger.Info("Exec command invoked")
			d.Logger.Debug("Getting ExecArgs `Project.GetExecArgs`")
			execArgs, err := d.Project.GetExecArgs()

			if err != nil {
				d.Logger.Error("Getting ExecArgs `Project.GetExecArgs` FAILED")
				return err
			}

			d.Logger.Debug("Calling `Backend.ExecutePipeline`")
			refID, res, err := d.Backend.ExecutePipeline(execArgs)

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
			execDetails, res, err := d.Backend.WaitForPipelineToFinish(refID, waitTimeout)
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

	return cmd
}
