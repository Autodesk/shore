package cleanup_command

import (
	"github.com/Autodesk/shore/pkg/command"
	"github.com/spf13/cobra"
)

// NewExecCommand - Using a Project, Renderer & Backend, executes a pipeline pipeline.
func NewCleanupCommand(d *command.Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Collection of cleanup pipeline related commands",
		Long: `Cleanup pipelines are procedures to remove resources that are created during a "main.pipeline.jsonnet" execution.
Use these subcommands when writing cleanup logic for the main pipeline.
When reversing, removing, and undoing previous operations that the main pipeline has created.
Cleanup pipelines are especially useful as cleanup procedure after testing.
		`,
	}

	cmd.AddCommand(
		NewRenderCommand(d),
		NewSaveCommand(d),
		NewExecCommand(d),
	)

	return cmd
}
