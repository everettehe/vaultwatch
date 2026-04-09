package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	cfgYAML := `
vault:
  address: "https://vault.example.com"
  token: "s.testtoken"
alerts:
  warn_before: 168h
  crit_before: 24h
  slack_webhook: "https://hooks.slack.com/test"
secrets:
  - path: secret/data/myapp/db
    description: "Database credentials"
`
	path := writeTempConfig(t, cfgYAML)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("unexpected address: %s", cfg.Vault.Address)
	}
	if cfg.Alerts.WarnBefore != 168*time.Hour {
		t.Errorf("unexpected warn_before: %v", cfg.Alerts.WarnBefore)
	}
	if len(cfg.Secrets) != 1 {
		t.Errorf("expected 1 secret, got %d", len(cfg.Secrets))
	}
}

func TestLoad_DefaultDurations(t *testing.T) {
	cfgYAML := `
vault:
  address: "https://vault.example.com"
secrets:
  - path: secret/data/myapp/key
`
	path := writeTempConfig(t, cfgYAML)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Alerts.WarnBefore != 7*24*time.Hour {
		t.Errorf("expected default warn_before 7d, got %v", cfg.Alerts.WarnBefore)
	}
	if cfg.Alerts.CritBefore != 24*time.Hour {
		t.Errorf("expected default crit_before 24h, got %v", cfg.Alerts.CritBefore)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	cfgYAML := `
secrets:
  - path: secret/data/myapp/key
`
	path := writeTempConfig(t, cfgYAML)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing vault.address")
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	cfgYAML := `
vault:
  address: "https://vault.example.com"
  token: "original"
secrets:
  - path: secret/data/myapp/key
`
	path := writeTempConfig(t, cfgYAML)
	t.Setenv("VAULT_TOKEN", "overridden")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Token != "overridden" {
		t.Errorf("expected token override, got %q", cfg.Vault.Token)
	}
}
