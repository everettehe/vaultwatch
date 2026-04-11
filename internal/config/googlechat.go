package config

// GoogleChatConfig holds configuration for the Google Chat notifier.
type GoogleChatConfig struct {
	// WebhookURL is the incoming webhook URL for the target Google Chat space.
	WebhookURL string `yaml:"webhook_url"`
}
