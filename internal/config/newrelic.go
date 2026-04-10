package config

// NewRelicConfig holds configuration for the New Relic notifier.
type NewRelicConfig struct {
	AccountID string `yaml:"account_id"`
	APIKey    string `yaml:"api_key"`
}
