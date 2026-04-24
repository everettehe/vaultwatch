package config

// AppConfigConfig holds configuration for the AWS AppConfig notifier.
type AppConfigConfig struct {
	// Application is the name or ID of the AppConfig application.
	Application string `mapstructure:"application"`

	// Environment is the name or ID of the AppConfig environment.
	Environment string `mapstructure:"environment"`

	// Profile is the name or ID of the AppConfig configuration profile.
	Profile string `mapstructure:"profile"`

	// Region is the AWS region. Defaults to us-east-1.
	Region string `mapstructure:"region"`
}
