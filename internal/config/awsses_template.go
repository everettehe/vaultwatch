package config

// SESTemplateConfig holds configuration for the AWS SES templated email notifier.
type SESTemplateConfig struct {
	From         string `yaml:"from"`
	To           string `yaml:"to"`
	TemplateName string `yaml:"template_name"`
	Region       string `yaml:"region"`
}
