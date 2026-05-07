package cli

import (
	"fmt"

	"github.com/nicholasgasior/envchain/internal/chain"
	"github.com/nicholasgasior/envchain/internal/envfile"
	"github.com/nicholasgasior/envchain/internal/output"
	"github.com/spf13/cobra"
)

func newSnapshotCmd() *cobra.Command {
	var configPath string
	var savePath string
	var compareWith string
	var format string

	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Save or compare environment snapshots",
		Long:  "Save the current merged environment to a snapshot file, or diff against a previously saved snapshot.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := chain.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			env, err := chain.Load(cfg.ToLayers())
			if err != nil {
				return fmt.Errorf("load chain: %w", err)
			}

			if compareWith != "" {
				snap, err := envfile.LoadSnapshot(compareWith)
				if err != nil {
					return fmt.Errorf("load snapshot: %w", err)
				}
				diff := envfile.DiffSnapshot(snap, env)
				if !envfile.HasChanges(diff) {
					fmt.Fprintln(cmd.OutOrStdout(), "No changes since snapshot.")
					return nil
				}
				return output.WriteDiff(cmd.OutOrStdout(), diff)
			}

			if savePath == "" {
				return fmt.Errorf("provide --save or --compare-with")
			}

			snap := envfile.NewSnapshot(configPath, env)
			if err := envfile.SaveSnapshot(savePath, snap); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved to %s\n", savePath)
			_ = format
			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", ".envchain.yml", "path to envchain config")
	cmd.Flags().StringVar(&savePath, "save", "", "path to save snapshot JSON")
	cmd.Flags().StringVar(&compareWith, "compare-with", "", "path to existing snapshot to diff against")
	cmd.Flags().StringVar(&format, "format", "dotenv", "output format for diff (dotenv|json)")

	return cmd
}
