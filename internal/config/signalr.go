package config

// SignalRConfig holds configuration for the Azure SignalR notifier.
type SignalRConfig struct {
	EndpointURL string `yaml:"endpoint_url"`
	AccessKey   string `yaml:"access_key"`
	Hub         string `yaml:"hub"`
}
