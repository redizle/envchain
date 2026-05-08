package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/chain"
	"github.com/envchain/envchain/internal/envfile"
)

func newSchemaCmd() *cobra.Command {
	var configPath string
	var schemaPath string

	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate merged env against a YAML schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := chain.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			env, err := chain.Load(cfg)
			if err != nil {
				return fmt.Errorf("load env: %w", err)
			}

			schema, err := envfile.LoadSchema(schemaPath)
			if err != nil {
				return fmt.Errorf("load schema: %w", err)
			}

			violations := envfile.ValidateSchema(env, schema)
			if len(violations) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "schema validation passed")
				return nil
			}

			fmt.Fprintf(cmd.ErrOrStderr(), "schema validation failed (%d issue(s)):\n", len(violations))
			for _, v := range violations {
				fmt.Fprintf(cmd.ErrOrStderr(), "  - %s\n", v.Error())
			}
			os.Exit(1)
			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", ".envchain.yaml", "path to envchain config")
	cmd.Flags().StringVarP(&schemaPath, "schema", "s", ".envschema.yaml", "path to schema file")
	return cmd
}
