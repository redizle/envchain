package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd builds the top-level envchain command with all subcommands attached.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envchain",
		Short: "Manage layered .env files with secret interpolation and diff support",
		Long: `envchain layers .env files across environments, resolves secrets from
AWS SSM, Vault, or Doppler, and provides diff, audit, lint, and template tools.`,
		SilenceUsage: true,
	}

	root.AddCommand(
		newPrintCmd(),
		newDiffCmd(),
		newRunCmd(),
		newValidateCmd(),
		newInitCmd(),
		newInterpolateCmd(),
		newAuditCmd(),
		newSnapshotCmd(),
		newSchemaCmd(),
		newTemplateCmd(),
	)

	return root
}
