package command

import (
	"log"

	"github.com/spf13/cobra"
)

// NewExecCommand - Using a Project, Renderer & Backend, executes a pipeline pipeline.
func NewExecCommand(d *Dependencies) *cobra.Command {

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
			_, res, err := d.Backend.ExecutePipeline(execArgs)

			if err != nil {
				return err
			}

			log.Println(res)
			return nil
		},
	}

	return cmd
}
