package config

// SESRawConfig holds configuration for the AWS SES raw email notifier.
type SESRawConfig struct {
	From   string `yaml:"from"`
	To     string `yaml:"to"`
	Region string `yaml:"region"`
}
