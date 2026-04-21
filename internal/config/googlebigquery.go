package config

// BigQueryConfig holds configuration for the Google BigQuery notifier.
type BigQueryConfig struct {
	ProjectID string `yaml:"project_id"`
	DatasetID string `yaml:"dataset_id"`
	TableID   string `yaml:"table_id"`
	APIKey    string `yaml:"api_key"`
}
