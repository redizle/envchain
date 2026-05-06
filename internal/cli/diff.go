package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envchain/internal/chain"
	"github.com/user/envchain/internal/output"
)

func newDiffCmd() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Show diff between env layers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := chain.LoadConfig(configFile)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			layers := cfg.ToLayers()
			if len(layers) < 2 {
				fmt.Fprintln(os.Stderr, "need at least 2 layers to diff")
				return nil
			}

			diffs, err := chain.DiffLayers(layers)
			if err != nil {
				return fmt.Errorf("diff layers: %w", err)
			}

			for i, d := range diffs {
				if len(d) == 0 {
					continue
				}
				fmt.Fprintf(os.Stdout, "--- layer %d → %d ---\n", i, i+1)
				if err := output.WriteDiff(os.Stdout, d); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", ".envchain.yaml", "path to envchain config file")

	return cmd
}
