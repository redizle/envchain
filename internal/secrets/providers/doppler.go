package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const dopplerBaseURL = "https://api.doppler.com/v3/configs/config/secret"

// DopplerClient is the interface for fetching secrets from Doppler.
type DopplerClient interface {
	GetSecret(ctx context.Context, token, project, config, name string) (string, error)
}

type httpDopplerClient struct {
	httpClient *http.Client
}

func (c *httpDopplerClient) GetSecret(ctx context.Context, token, project, config, name string) (string, error) {
	url := fmt.Sprintf("%s?project=%s&config=%s&name=%s", dopplerBaseURL, project, config, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(token, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("secret %q not found in doppler", name)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("doppler returned status %d", resp.StatusCode)
	}

	var result struct {
		Secret struct {
			Raw string `json:"raw"`
		} `json:"secret"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding doppler response: %w", err)
	}
	return result.Secret.Raw, nil
}

// DopplerProvider resolves secrets from Doppler using the format:
// doppler:PROJECT/CONFIG/SECRET_NAME
type DopplerProvider struct {
	token  string
	client DopplerClient
}

// NewDopplerProvider creates a provider using the real HTTP client.
func NewDopplerProvider(token string) *DopplerProvider {
	return &DopplerProvider{token: token, client: &httpDopplerClient{httpClient: &http.Client{}}}
}

// NewDopplerProviderWithClient creates a provider with a custom client (for testing).
func NewDopplerProviderWithClient(token string, client DopplerClient) *DopplerProvider {
	return &DopplerProvider{token: token, client: client}
}

func (p *DopplerProvider) Name() string { return "doppler" }

func (p *DopplerProvider) Resolve(ctx context.Context, ref string) (string, error) {
	parts := strings.SplitN(ref, "/", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("doppler ref must be PROJECT/CONFIG/SECRET_NAME, got %q", ref)
	}
	return p.client.GetSecret(ctx, p.token, parts[0], parts[1], parts[2])
}
