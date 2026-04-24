package config

// AWSIoTConfig holds configuration for the AWS IoT notifier.
type AWSIoTConfig struct {
	// Topic is the MQTT topic to publish alerts to (e.g. "vaultwatch/alerts").
	Topic string `yaml:"topic"`
	// Region is the AWS region where the IoT endpoint resides.
	Region string `yaml:"region"`
}
