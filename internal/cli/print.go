package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envchain/internal/chain"
	"github.com/user/envchain/internal/output"
	"github.com/user/envchain/internal/secrets"
	"github.com/user/envchain/internal/secrets/providers"
)

func newPrintCmd() *cobra.Command {
	var (
		configFile string
		format     string
		resolve    bool
	)

	cmd := &cobra.Command{
		Use:   "print",
		Short: "Print merged environment variables",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := chain.LoadConfig(configFile)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			env, err := chain.Load(cfg.ToLayers())
			if err != nil {
				return fmt.Errorf("load layers: %w", err)
			}

			if resolve {
				env, err = resolveSecrets(env)
				if err != nil {
					return err
				}
			}

			return output.WriteEnv(os.Stdout, env, format)
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", ".envchain.yaml", "path to envchain config file")
	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "output format: dotenv, export, json")
	cmd.Flags().BoolVar(&resolve, "resolve", false, "resolve secret references")

	return cmd
}

// resolveSecrets resolves any secret references in the given environment map
// using the default provider registry.
func resolveSecrets(env map[string]string) (map[string]string, error) {
	reg := providers.NewRegistry()
	r := secrets.NewResolver(reg)
	resolved, err := r.Resolve(env)
	if err != nil {
		return nil, fmt.Errorf("resolve secrets: %w", err)
	}
	return resolved, nil
}
