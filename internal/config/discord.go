package config

// DiscordConfig holds configuration for the Discord notifier.
type DiscordConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}
