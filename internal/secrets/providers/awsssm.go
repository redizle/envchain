// Package providers implements secret provider backends for envchain.
package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SSMProvider resolves secrets from AWS Systems Manager Parameter Store.
type SSMProvider struct {
	client *ssm.Client
}

// NewSSMProvider creates a new SSMProvider using the default AWS config.
func NewSSMProvider(ctx context.Context) (*SSMProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("awsssm: load config: %w", err)
	}
	return &SSMProvider{client: ssm.NewFromConfig(cfg)}, nil
}

// NewSSMProviderWithClient creates an SSMProvider with a pre-configured client.
// Useful for testing with a mock client.
func NewSSMProviderWithClient(client *ssm.Client) *SSMProvider {
	return &SSMProvider{client: client}
}

// Name returns the provider identifier used in secret references (e.g. awsssm://my-param).
func (p *SSMProvider) Name() string {
	return "awsssm"
}

// Resolve fetches a parameter from SSM Parameter Store.
// The ref should be the parameter name or path, e.g. "/myapp/prod/db_password".
func (p *SSMProvider) Resolve(ctx context.Context, ref string) (string, error) {
	ref = strings.TrimPrefix(ref, "/")
	paramName := "/" + ref

	out, err := p.client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("awsssm: get parameter %q: %w", paramName, err)
	}
	if out.Parameter == nil || out.Parameter.Value == nil {
		return "", fmt.Errorf("awsssm: parameter %q returned nil value", paramName)
	}
	return aws.ToString(out.Parameter.Value), nil
}
