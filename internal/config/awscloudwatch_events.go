package config

// CloudWatchEventsConfig holds configuration for the AWS CloudWatch Events notifier.
type CloudWatchEventsConfig struct {
	EventBus   string `yaml:"event_bus"`
	Source     string `yaml:"source"`
	DetailType string `yaml:"detail_type"`
	Region     string `yaml:"region"`
}
