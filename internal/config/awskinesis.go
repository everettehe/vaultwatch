package config

// KinesisConfig holds configuration for the AWS Kinesis Data Stream notifier.
type KinesisConfig struct {
	// StreamName is the name of the Kinesis Data Stream to publish events to.
	StreamName string `yaml:"stream_name"`

	// PartitionKey is used to group related records in the stream.
	// Defaults to "vaultwatch" if not set.
	PartitionKey string `yaml:"partition_key"`
}
