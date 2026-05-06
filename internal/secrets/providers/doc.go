// Package providers contains implementations of secret provider backends
// for envchain. Each provider satisfies the Provider interface and can be
// registered with a Registry for dispatch during secret resolution.
//
// Built-in providers:
//
//   - awsssm: resolves secrets from AWS Systems Manager Parameter Store
//   - vault:  resolves secrets from HashiCorp Vault KV engine
//
// # Adding a new provider
//
// Implement the Provider interface:
//
//	type Provider interface {
//	    Name() string
//	    Resolve(ctx context.Context, ref string) (string, error)
//	}
//
// Then register it with a Registry before resolving secrets:
//
//	reg := providers.NewRegistry()
//	reg.Register(providers.NewSSMProvider())
//	reg.Register(providers.NewVaultProvider(url, token))
package providers
