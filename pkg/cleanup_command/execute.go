package cleanup_command

import (
	"github.com/Autodeskshore/pkg/command"
	"github.com/spf13/cobra"
)

// NewExecCommand - Using a Project, Renderer & Backend, executes a pipeline pipeline.
func NewExecCommand(d *command.Dependencies) *cobra.Command {
	cmd := command.NewExecCommand(d, "cleanup/exec")
	cmd.Use = "exec"
	cmd.Short = "execute the cleanup pipeline"
	cmd.Long = "exec"

	return cmd
}
