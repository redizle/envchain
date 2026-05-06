package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd builds the root cobra command with all subcommands attached.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envchain",
		Short: "Manage layered .env files with secret interpolation and diff support",
		Long: `envchain loads, merges, and resolves layered .env files.
Secrets can be pulled from AWS SSM or Vault using ref syntax: secret://provider/path`,
		SilenceUsage: true,
	}

	root.AddCommand(newPrintCmd())
	root.AddCommand(newDiffCmd())

	return root
}
