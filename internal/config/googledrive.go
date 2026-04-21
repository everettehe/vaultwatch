package config

// GoogleDriveConfig holds configuration for the Google Drive (Sheets) notifier.
type GoogleDriveConfig struct {
	SpreadsheetID string `yaml:"spreadsheet_id"`
	SheetName     string `yaml:"sheet_name"`
	APIKey        string `yaml:"api_key"`
}
