package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/chain"
	"github.com/envchain/envchain/internal/envfile"
	"github.com/envchain/envchain/internal/output"
)

func newTemplateCmd() *cobra.Command {
	var (
		configPath string
		env        string
		missingKey string
		format     string
	)

	cmd := &cobra.Command{
		Use:   "template",
		Short: "Render template expressions inside env values",
		Long: `Loads the env chain, then renders any Go template expressions
found in values (e.g. {{ .OTHER_KEY }} or {{ default "x" .KEY }}).
The result is printed in the requested format.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := chain.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			merged, err := chain.Load(cfg, env)
			if err != nil {
				return fmt.Errorf("load chain: %w", err)
			}

			opts := envfile.TemplateOptions{MissingKey: missingKey}
			rendered, err := envfile.RenderTemplateMap(merged, opts)
			if err != nil {
				return fmt.Errorf("render templates: %w", err)
			}

			return output.WriteEnv(os.Stdout, rendered, output.Format(format))
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", ".envchain.yaml", "path to envchain config")
	cmd.Flags().StringVarP(&env, "env", "e", "development", "environment name")
	cmd.Flags().StringVar(&missingKey, "missing-key", "error", "behaviour on missing key: error, zero, default")
	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "output format: dotenv, export, json, docker")

	return cmd
}
