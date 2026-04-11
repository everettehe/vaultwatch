package config

// LarkConfig holds configuration for Lark (Feishu) webhook notifications.
type LarkConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}
