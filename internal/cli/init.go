package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var defaultConfig = `# envchain configuration
layers:
  - base.env
  - local.env
`

var defaultBaseEnv = `# Base environment variables
APP_ENV=development
APP_PORT=8080
`

var defaultLocalEnv = `# Local overrides (do not commit)
# APP_SECRET=ref:ssm:/myapp/secret
`

func newInitCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new envchain project in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			files := map[string]string{
				"envchain.yaml": defaultConfig,
				"base.env":      defaultBaseEnv,
				"local.env":     defaultLocalEnv,
			}

			for name, content := range files {
				path := filepath.Join(".", name)
				if !force {
					if _, err := os.Stat(path); err == nil {
						fmt.Fprintf(cmd.OutOrStdout(), "skipping %s (already exists, use --force to overwrite)\n", name)
						continue
					}
				}
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					return fmt.Errorf("writing %s: %w", name, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "created %s\n", name)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "\nenvchain project initialized. Edit envchain.yaml to configure your layers.")
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing files")
	return cmd
}
