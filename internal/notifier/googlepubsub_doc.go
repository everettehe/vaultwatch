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
package notifier
