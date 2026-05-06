package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain/internal/chain"
	"github.com/yourorg/envchain/internal/envfile"
)

func newValidateCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate merged env variables for key/value correctness",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := chain.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			merged, err := chain.Load(cfg)
			if err != nil {
				return fmt.Errorf("loading chain: %w", err)
			}

			if err := envfile.Validate(merged); err != nil {
				return fmt.Errorf("validation errors:\n%w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "OK: all variables are valid")
			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", ".envchain.yaml", "path to envchain config file")
	return cmd
}
