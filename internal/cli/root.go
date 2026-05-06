package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd builds the root cobra command with all subcommands attached.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envchain",
		Short: "Manage layered .env files with secret interpolation",
		Long: `envchain lets you compose .env files across environments,
interpolate secrets from providers like AWS SSM, Vault, and Doppler,
and diff or audit your configuration.`,
		SilenceUsage: true,
	}

	root.AddCommand(newPrintCmd())
	root.AddCommand(newDiffCmd())
	root.AddCommand(newRunCmd())
	root.AddCommand(newValidateCmd())
	root.AddCommand(newInitCmd())
	root.AddCommand(newInterpolateCmd())
	root.AddCommand(newAuditCmd())

	return root
}
