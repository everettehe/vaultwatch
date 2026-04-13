package config

// LineNotifyConfig holds configuration for the LINE Notify notifier.
type LineNotifyConfig struct {
	// Token is the LINE Notify personal access token.
	Token string `yaml:"token"`
}
