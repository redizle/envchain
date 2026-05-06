package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain/internal/chain"
	"github.com/yourorg/envchain/internal/envfile"
)

func newAuditCmd() *cobra.Command {
	var configPath string
	var env string

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit merged env for security and quality issues",
		Long: `Audit inspects the merged environment for common problems such as
plain-text secrets, literal newlines in values, and unusually long values.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := chain.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			layers := cfg.ToLayers(env)
			merged, err := chain.Load(layers)
			if err != nil {
				return fmt.Errorf("loading layers: %w", err)
			}

			result := envfile.Audit(merged)
			if !result.HasIssues() {
				fmt.Fprintln(cmd.OutOrStdout(), "No issues found.")
				return nil
			}

			fmt.Fprint(cmd.OutOrStdout(), result.String())

			if len(result.Errors) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", ".envchain.yaml", "Path to envchain config file")
	cmd.Flags().StringVarP(&env, "env", "e", "", "Target environment (e.g. production)")

	return cmd
}
