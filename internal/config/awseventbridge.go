package config

// EventBridgeConfig holds configuration for the AWS EventBridge notifier.
type EventBridgeConfig struct {
	// EventBus is the event bus name or ARN. Defaults to "default".
	EventBus string `yaml:"event_bus"`

	// Source is the event source string. Defaults to "vaultwatch".
	Source string `yaml:"source"`

	// DetailType is the EventBridge detail-type field. Defaults to "VaultSecretExpiry".
	DetailType string `yaml:"detail_type"`
}
