package config

// AWSConfigConfig holds configuration for the AWS Config notifier.
type AWSConfigConfig struct {
	ResultToken string `yaml:"result_token"`
	Region      string `yaml:"region"`
}
