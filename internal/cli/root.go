package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd builds the top-level cobra command with all subcommands attached.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envchain",
		Short: "Manage layered .env files with secret interpolation",
		SilenceUsage: true,
	}

	root.AddCommand(newPrintCmd())
	root.AddCommand(newDiffCmd())
	root.AddCommand(newRunCmd())
	root.AddCommand(newValidateCmd())

	return root
}
