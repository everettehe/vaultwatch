package config

// GoogleTaskConfig holds configuration for the Google Task notifier.
type GoogleTaskConfig struct {
	WebhookURL string `yaml:"webhook_url"`
	Tasklist   string `yaml:"tasklist"`
}
