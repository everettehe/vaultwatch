package config

// GoogleStorageConfig holds settings for the Google Cloud Storage notifier.
type GoogleStorageConfig struct {
	Bucket string `mapstructure:"bucket"`
	APIKey string `mapstructure:"api_key"`
}
