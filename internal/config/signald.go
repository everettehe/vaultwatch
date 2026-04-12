package config

// SignaldConfig holds configuration for the Signald notifier.
type SignaldConfig struct {
	BaseURL    string   `yaml:"base_url"`
	Sender     string   `yaml:"sender"`
	Recipients []string `yaml:"recipients"`
}
