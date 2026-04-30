package config

// AWSConnectConfig holds configuration for the Amazon Connect notifier.
type AWSConnectConfig struct {
	InstanceID  string `yaml:"instance_id"`
	ContactFlow string `yaml:"contact_flow"`
	QueueID     string `yaml:"queue_id"`
	Region      string `yaml:"region"`
}
