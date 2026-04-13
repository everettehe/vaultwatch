// Package notifier provides alert delivery implementations for VaultWatch.
//
// # Google Cloud Pub/Sub Notifier
//
// The GooglePubSubNotifier publishes secret expiration alerts to a Google Cloud
// Pub/Sub topic. It requires a GCP project ID and topic ID. Authentication is
// handled via Application Default Credentials (ADC).
//
// # Configuration
//
//	[notifiers.googlepubsub]
//	project_id = "my-gcp-project"
//	topic_id   = "vault-alerts"
//
// # Authentication
//
// Set GOOGLE_APPLICATION_CREDENTIALS or use Workload Identity when running
// on GKE/Cloud Run.
//
// # Message Format
//
// Each published message contains a JSON payload with the following fields:
//
//	{
//	  "secret_path": "secret/myapp/db-password",
//	  "expires_at":  "2024-06-01T00:00:00Z",
//	  "severity":    "critical"
//	}
//
// Messages also include a "severity" attribute on the Pub/Sub message itself
// to allow subscription filtering (e.g. filter by severity="critical").
package notifier
