// Package notifier provides notification backends for VaultWatch.
//
// # Google BigQuery Notifier
//
// The BigQueryNotifier streams secret expiration events into a Google BigQuery
// table using the BigQuery REST streaming insert API.
//
// # Configuration
//
//	[notifiers.bigquery]
//	project_id = "my-gcp-project"
//	dataset_id = "vaultwatch"
//	table_id   = "secret_expirations"
//	api_key    = "AIza..."
//
// # Table Schema
//
// The notifier writes rows with the following fields:
//
//	- path        STRING   — vault secret path
//	- days_left   INT64    — days remaining until expiration
//	- expires_at  STRING   — RFC3339 expiration timestamp
//	- is_expired  BOOL     — whether the secret is already expired
//	- notified_at STRING   — RFC3339 timestamp of the notification
package notifier
