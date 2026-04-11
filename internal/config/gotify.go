package config

// GotifyConfig holds configuration for the Gotify notifier.
type GotifyConfig struct {
	ServerURL string `yaml:"server_url"`
	Token     string `yaml:"token"`
	Priority  int    `yaml:"priority"`
}
