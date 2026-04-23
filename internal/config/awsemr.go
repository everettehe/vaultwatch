package config

// EMRConfig holds configuration for the AWS EMR notifier.
type EMRConfig struct {
	ClusterID string `yaml:"cluster_id"`
	Region    string `yaml:"region"`
}
