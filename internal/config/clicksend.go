package config

// ClickSendConfig holds configuration for the ClickSend SMS notifier.
type ClickSendConfig struct {
	Username string `yaml:"username"`
	APIKey   string `yaml:"api_key"`
	To       string `yaml:"to"`
}
