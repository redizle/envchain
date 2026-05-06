package cli

import (
	"os"

	"github.com/user/envchain/internal/secrets/providers"
)

// buildRegistry constructs a provider registry from environment-based
// configuration. Providers are only registered when their required
// environment variables are present.
func buildRegistry() *providers.Registry {
	reg := providers.NewRegistry()

	// AWS SSM — requires AWS credentials to be configured externally;
	// we register it unconditionally since the SDK handles auth.
	reg.Register(providers.NewSSMProvider())

	// Vault — requires VAULT_ADDR and VAULT_TOKEN.
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		token := os.Getenv("VAULT_TOKEN")
		reg.Register(providers.NewVaultProvider(addr, token))
	}

	// Doppler — requires DOPPLER_TOKEN.
	if tok := os.Getenv("DOPPLER_TOKEN"); tok != "" {
		reg.Register(providers.NewDopplerProvider(tok))
	}

	return reg
}
