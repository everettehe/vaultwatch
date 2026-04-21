// Package notifier provides alert delivery integrations for vaultwatch.
//
// # Google Drive (Sheets) Notifier
//
// GoogleDriveNotifier appends a row to a Google Sheets spreadsheet whenever
// a secret is expiring or has expired. Each row contains:
//
//   - Timestamp (RFC3339)
//   - Secret path
//   - Days until expiration
//   - Alert message body
//
// # Configuration
//
//	notifiers:
//	  googledrive:
//	    spreadsheet_id: "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms"
//	    sheet_name: "VaultAlerts"   # optional, defaults to "Alerts"
//	    api_key: "AIza..."
//
// The API key must have write access to the target spreadsheet.
package notifier
