package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/chain"
	"github.com/envchain/envchain/internal/secrets"
	"github.com/envchain/envchain/internal/secrets/providers"
)

func newRunCmd() *cobra.Command {
	var (
		configFile string
		env        string
	)

	cmd := &cobra.Command{
		Use:   "run [flags] -- <command> [args...]",
		Short: "Run a command with resolved env vars injected",
		Example: `  envchain run --env production -- printenv
  envchain run --config .envchain.yaml -- node server.js`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := chain.LoadConfig(configFile)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			layers := cfg.ToLayers(env)
			merged, err := chain.Load(layers)
			if err != nil {
				return fmt.Errorf("loading env chain: %w", err)
			}

			reg := providers.NewRegistry()
			resolver := secrets.NewResolver(reg)
			resolved, err := resolver.Resolve(merged)
			if err != nil {
				return fmt.Errorf("resolving secrets: %w", err)
			}

			envVars := os.Environ()
			for k, v := range resolved {
				envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
			}

			binary, err := exec.LookPath(args[0])
			if err != nil {
				return fmt.Errorf("command not found: %s", args[0])
			}

			proc := exec.Command(binary, args[1:]...)
			proc.Env = envVars
			proc.Stdin = os.Stdin
			proc.Stdout = os.Stdout
			proc.Stderr = os.Stderr

			if err := proc.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					os.Exit(exitErr.ExitCode())
				}
				return fmt.Errorf("running command %q: %w", strings.Join(args, " "), err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&configFile, "config", ".envchain.yaml", "path to config file")
	cmd.Flags().StringVar(&env, "env", "", "target environment (e.g. production, staging)")

	return cmd
}
