package config

// SQSConfig holds configuration for the AWS SQS notifier.
type SQSConfig struct {
	// QueueURL is the full URL of the SQS queue to send messages to.
	QueueURL string `mapstructure:"queue_url"`
}
