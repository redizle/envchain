package cli

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/chain"
	"github.com/envchain/envchain/internal/envfile"
)

func newEncryptCmd() *cobra.Command {
	var (
		cfgPath string
		keyHex  string
		decrypt bool
	)

	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt or decrypt values in the resolved env",
		Long: `Encrypt or decrypt all values in the merged env using AES-256-GCM.
Provide a 32-byte key as a 64-character hex string via --key or ENVCHAIN_KEY.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if keyHex == "" {
				keyHex = os.Getenv("ENVCHAIN_KEY")
			}
			if len(keyHex) != 64 {
				return fmt.Errorf("--key must be a 64-character hex string (32 bytes)")
			}

			key, err := hex.DecodeString(keyHex)
			if err != nil {
				return fmt.Errorf("invalid hex key: %w", err)
			}

			cfg, err := chain.LoadConfig(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			env, err := chain.Load(cfg)
			if err != nil {
				return fmt.Errorf("load env: %w", err)
			}

			var result map[string]string
			if decrypt {
				result, err = envfile.DecryptMap(env, key)
				if err != nil {
					return fmt.Errorf("decrypt: %w", err)
				}
			} else {
				result, err = envfile.EncryptMap(env, key)
				if err != nil {
					return fmt.Errorf("encrypt: %w", err)
				}
			}

			for k, v := range result {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&cfgPath, "config", "c", ".envchain.yaml", "path to envchain config")
	cmd.Flags().StringVar(&keyHex, "key", "", "64-char hex AES-256 key (or set ENVCHAIN_KEY)")
	cmd.Flags().BoolVar(&decrypt, "decrypt", false, "decrypt values instead of encrypting")

	return cmd
}
