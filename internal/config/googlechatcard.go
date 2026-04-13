package config

// GoogleChatCardConfig holds configuration for the Google Chat card notifier.
type GoogleChatCardConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}
