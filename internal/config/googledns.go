package config

// GoogleDNSConfig holds configuration for the Google DNS webhook notifier.
type GoogleDNSConfig struct {
	WebhookURL string `yaml:"webhook_url"`
	Project    string `yaml:"project"`
}
