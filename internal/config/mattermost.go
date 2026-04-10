package config

// MattermostConfig holds configuration for the Mattermost notifier.
type MattermostConfig struct {
	// WebhookURL is the incoming webhook URL created in Mattermost.
	WebhookURL string `yaml:"webhook_url"`

	// Channel overrides the default channel configured in the webhook (optional).
	Channel string `yaml:"channel"`

	// Username overrides the display name for posted messages (optional).
	Username string `yaml:"username"`
}
