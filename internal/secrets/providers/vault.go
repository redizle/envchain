package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"strings"
)

// VaultClient is an interface for fetching secrets from Vault.
type VaultClient interface {
	GetSecret(ctx context.Context, path string) (string, error)
}

// httpVaultClient is the real Vault HTTP client.
type httpVaultClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func (c *httpVaultClient) GetSecret(ctx context.Context, path string) (string, error) {
	url := fmt.Sprintf("%s/v1/%s", strings.TrimRight(c.baseURL, "/"), path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vault: unexpected status %d for path %q", resp.StatusCode, path)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Data struct {
			Value string `json:"value"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("vault: failed to parse response: %w", err)
	}
	return result.Data.Value, nil
}

// VaultProvider resolves secrets from HashiCorp Vault.
type VaultProvider struct {
	client VaultClient
}

// NewVaultProvider creates a VaultProvider using the real HTTP client.
func NewVaultProvider(baseURL, token string) *VaultProvider {
	return &VaultProvider{
		client: &httpVaultClient{
			baseURL: baseURL,
			token:   token,
			client:  &http.Client{},
		},
	}
}

// NewVaultProviderWithClient creates a VaultProvider with a custom client (for testing).
func NewVaultProviderWithClient(client VaultClient) *VaultProvider {
	return &VaultProvider{client: client}
}

// Name returns the provider identifier used in secret refs.
func (p *VaultProvider) Name() string { return "vault" }

// Resolve fetches the secret value at the given path from Vault.
func (p *VaultProvider) Resolve(ctx context.Context, ref string) (string, error) {
	return p.client.GetSecret(ctx, ref)
}
