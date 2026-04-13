package config

// GoogleSheetsConfig holds configuration for the Google Sheets notifier.
type GoogleSheetsConfig struct {
	// WebAppURL is the Google Apps Script Web App deployment URL.
	WebAppURL string `mapstructure:"web_app_url"`
	// SheetName is the optional name of the target sheet tab.
	SheetName string `mapstructure:"sheet_name"`
}
