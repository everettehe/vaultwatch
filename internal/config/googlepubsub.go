package config

// GooglePubSubConfig holds configuration for the Google Cloud Pub/Sub notifier.
type GooglePubSubConfig struct {
	// ProjectID is the GCP project that owns the topic.
	ProjectID string `yaml:"project_id"`

	// TopicID is the Pub/Sub topic name (without the full resource path).
	TopicID string `yaml:"topic_id"`
}
