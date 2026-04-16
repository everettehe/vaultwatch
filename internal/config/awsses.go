package config

// SESConfig holds configuration for the AWS SES email notifier.
type SESConfig struct {
	From   string `yaml:"from"`
	To     string `yaml:"to"`
	Region string `yaml:"region"`
}
