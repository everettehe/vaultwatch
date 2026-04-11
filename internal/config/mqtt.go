package config

// MQTTConfig holds configuration for the MQTT notifier.
type MQTTConfig struct {
	// BrokerURL is the MQTT broker address, e.g. "tcp://localhost:1883".
	BrokerURL string `yaml:"broker_url"`

	// Topic is the MQTT topic to publish alerts to.
	Topic string `yaml:"topic"`

	// ClientID is an optional MQTT client identifier.
	// Defaults to a generated value if empty.
	ClientID string `yaml:"client_id"`

	// QoS is the MQTT quality-of-service level (0, 1, or 2). Defaults to 1.
	QoS byte `yaml:"qos"`
}
