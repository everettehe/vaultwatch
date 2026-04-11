package config

// NtfyConfig holds configuration for the ntfy notifier.
type NtfyConfig struct {
	// ServerURL is the base URL of the ntfy server.
	// Defaults to "https://ntfy.sh" if not set.
	ServerURL string `yaml:"server_url"`

	// Topic is the ntfy topic to publish alerts to.
	// Required.
	Topic string `yaml:"topic"`
}
