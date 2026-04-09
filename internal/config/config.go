package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level vaultwatch configuration.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Alerts  AlertsConfig  `yaml:"alerts"`
	Secrets []SecretWatch `yaml:"secrets"`
}

// VaultConfig contains Vault connection settings.
type VaultConfig struct {
	Address   string `yaml:"address"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
}

// AlertsConfig defines when and how to send alerts.
type AlertsConfig struct {
	WarnBefore  time.Duration `yaml:"warn_before"`
	CritBefore  time.Duration `yaml:"crit_before"`
	SlackWebhook string       `yaml:"slack_webhook"`
	Email        string       `yaml:"email"`
}

// SecretWatch describes a single secret to monitor.
type SecretWatch struct {
	Path        string `yaml:"path"`
	Description string `yaml:"description"`
}

// Load reads and parses the config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Allow token override via environment variable.
	if tok := os.Getenv("VAULT_TOKEN"); tok != "" {
		cfg.Vault.Token = tok
	}
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		cfg.Vault.Address = addr
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if len(c.Secrets) == 0 {
		return fmt.Errorf("at least one secret must be defined under secrets")
	}
	if c.Alerts.WarnBefore == 0 {
		c.Alerts.WarnBefore = 7 * 24 * time.Hour // default 7 days
	}
	if c.Alerts.CritBefore == 0 {
		c.Alerts.CritBefore = 24 * time.Hour // default 1 day
	}
	return nil
}
