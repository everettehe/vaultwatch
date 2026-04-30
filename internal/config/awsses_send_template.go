package config

// SESSendTemplateConfig holds configuration for the AWS SES templated email notifier.
type SESSendTemplateConfig struct {
	// From is the verified sender address.
	From string `yaml:"from"`
	// To is the recipient address.
	To string `yaml:"to"`
	// Template is the name of the SES template to use.
	Template string `yaml:"template"`
	// Region is the AWS region where the template is registered.
	Region string `yaml:"region"`
}
