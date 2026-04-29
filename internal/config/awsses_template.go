package config

// SESTemplateConfig holds configuration for the AWS SES Template notifier.
type SESTemplateConfig struct {
	// From is the sender email address.
	From string `yaml:"from"`
	// To is the recipient email address.
	To string `yaml:"to"`
	// Template is the name of the SES email template to use.
	Template string `yaml:"template"`
	// Region is the AWS region where SES is configured.
	Region string `yaml:"region"`
}
