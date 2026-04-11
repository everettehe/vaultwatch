package config

// StatusPageConfig holds Atlassian Statuspage notifier configuration.
type StatusPageConfig struct {
	APIKey string `yaml:"api_key"`
	PageID string `yaml:"page_id"`
}
