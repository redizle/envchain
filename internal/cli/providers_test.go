package cli

import (
	"os"
	"testing"
)

func TestBuildRegistry_AlwaysHasSSM(t *testing.T) {
	reg := buildRegistry()
	_, err := reg.Get("ssm")
	if err != nil {
		t.Errorf("expected ssm provider to always be registered: %v", err)
	}
}

func TestBuildRegistry_VaultRegisteredWhenEnvSet(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://localhost:8200")
	t.Setenv("VAULT_TOKEN", "root")

	reg := buildRegistry()
	_, err := reg.Get("vault")
	if err != nil {
		t.Errorf("expected vault provider when VAULT_ADDR is set: %v", err)
	}
}

func TestBuildRegistry_VaultAbsentWhenEnvMissing(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")

	reg := buildRegistry()
	_, err := reg.Get("vault")
	if err == nil {
		t.Error("expected vault provider to be absent without VAULT_ADDR")
	}
}

func TestBuildRegistry_DopplerRegisteredWhenEnvSet(t *testing.T) {
	t.Setenv("DOPPLER_TOKEN", "dp.st.abc123")

	reg := buildRegistry()
	_, err := reg.Get("doppler")
	if err != nil {
		t.Errorf("expected doppler provider when DOPPLER_TOKEN is set: %v", err)
	}
}

func TestBuildRegistry_DopplerAbsentWhenEnvMissing(t *testing.T) {
	os.Unsetenv("DOPPLER_TOKEN")

	reg := buildRegistry()
	_, err := reg.Get("doppler")
	if err == nil {
		t.Error("expected doppler provider to be absent without DOPPLER_TOKEN")
	}
}
