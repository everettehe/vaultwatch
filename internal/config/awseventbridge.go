package config

// AWSEventBridgeConfig holds configuration for the AWS EventBridge notifier.
type AWSEventBridgeConfig struct {
	// EventBus is the name or ARN of the EventBridge event bus.
	EventBus string `mapstructure:"event_bus"`

	// Source is the event source string (default: "vaultwatch").
	Source string `mapstructure:"source"`

	// Region is the AWS region. If empty, the SDK default chain is used.
	Region string `mapstructure:"region"`
}
