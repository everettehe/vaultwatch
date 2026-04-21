package config

// SQSFIFOConfig holds configuration for the AWS SQS FIFO notifier.
type SQSFIFOConfig struct {
	// QueueURL is the full URL of the SQS FIFO queue.
	QueueURL string `yaml:"queue_url"`
	// MessageGroupID is the message group ID used for FIFO ordering.
	MessageGroupID string `yaml:"message_group_id"`
}
