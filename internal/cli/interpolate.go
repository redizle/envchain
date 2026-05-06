package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/chain"
	"github.com/envchain/envchain/internal/envfile"
	"github.com/envchain/envchain/internal/output"
)

func newInterpolateCmd() *cobra.Command {
	var configPath string
	var format string
	var strict bool

	cmd := &cobra.Command{
		Use:   "interpolate",
		Short: "Resolve ${VAR} references within the merged env",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := chain.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			env, err := chain.Load(cfg)
			if err != nil {
				return fmt.Errorf("load chain: %w", err)
			}

			resolved, errs := envfile.Interpolate(env)
			if len(errs) > 0 {
				for _, e := range errs {
					fmt.Fprintf(os.Stderr, "warning: %v\n", e)
				}
				if strict {
					return fmt.Errorf("%d interpolation error(s)", len(errs))
				}
			}

			return output.WriteEnv(cmd.OutOrStdout(), resolved, format)
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", ".envchain.yaml", "path to envchain config")
	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "output format: dotenv, export, json")
	cmd.Flags().BoolVar(&strict, "strict", false, "exit non-zero if any references are unresolved")

	return cmd
}
