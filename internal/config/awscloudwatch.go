package config

// CloudWatchConfig holds configuration for the AWS CloudWatch notifier.
type CloudWatchConfig struct {
	// Namespace is the CloudWatch metric namespace. Defaults to "VaultWatch".
	Namespace string `yaml:"namespace"`
}
