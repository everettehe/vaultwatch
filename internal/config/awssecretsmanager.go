package config

// AWSSecretsManagerConfig holds configuration for the AWS Secrets Manager notifier.
type AWSSecretsManagerConfig struct {
	// SecretID is the ARN or name of the secret to write expiration events to.
	SecretID string `mapstructure:"secret_id"`
	// Region is the AWS region. Falls back to AWS_REGION env var if empty.
	Region string `mapstructure:"region"`
}
