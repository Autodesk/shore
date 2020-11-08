package controller

import (
	"log"

	"github.com/spf13/cobra"
)

// NewExecCommand - Using a Project, Renderer & Backend, executes a pipeline pipeline.
func NewExecCommand(d *Dependencies) *cobra.Command {
	var WithSave bool

	cmd := &cobra.Command{
		Use:   "exec",
		Short: "executes the pipeline",
		Long:  "Executes the selected pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			if WithSave {
				pipeline, err := Render(d)

				if err != nil {
					return err
				}

				_, err = d.Backend.SavePipeline(pipeline)

				if err != nil {
					return err
				}
			}

			execArgs, err := d.Project.GetExecArgs()

			if err != nil {
				return err
			}

			res, err := d.Backend.ExecutePipeline(execArgs)

			if err != nil {
				return err
			}

			log.Println(res)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&WithSave, "save", "s", false, "Render & Save the pipeline before executing it")

	return cmd
}
