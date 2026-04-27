package config

// EKSConfig holds configuration for the AWS EKS notifier.
type EKSConfig struct {
	Cluster string `yaml:"cluster"`
	Addon   string `yaml:"addon"`
	Region  string `yaml:"region"`
}
