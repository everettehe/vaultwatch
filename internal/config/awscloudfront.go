package config

// CloudFrontConfig holds configuration for the AWS CloudFront notifier.
type CloudFrontConfig struct {
	DistributionID string   `yaml:"distribution_id"`
	Region         string   `yaml:"region"`
	Paths          []string `yaml:"paths"`
}
