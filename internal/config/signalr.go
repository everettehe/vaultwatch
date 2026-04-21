package config

import "errors"

// SignalRConfig holds configuration for the Azure SignalR notifier.
type SignalRConfig struct {
	EndpointURL string `yaml:"endpoint_url"`
	AccessKey   string `yaml:"access_key"`
	Hub         string `yaml:"hub"`
}

// Validate checks that all required SignalRConfig fields are set.
func (c *SignalRConfig) Validate() error {
	if c.EndpointURL == "" {
		return errors.New("signalr: endpoint_url is required")
	}
	if c.AccessKey == "" {
		return errors.New("signalr: access_key is required")
	}
	if c.Hub == "" {
		return errors.New("signalr: hub is required")
	}
	return nil
}
