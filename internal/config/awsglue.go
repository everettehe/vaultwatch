package config

// GlueConfig holds configuration for the AWS Glue notifier.
type GlueConfig struct {
	// JobName is the name of the Glue job to trigger on secret expiry events.
	JobName string `yaml:"job_name" mapstructure:"job_name"`

	// Region is the AWS region where the Glue job is deployed.
	Region string `yaml:"region" mapstructure:"region"`
}
