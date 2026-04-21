// Package notifier provides alert delivery integrations for vaultwatch.
//
// # Google DNS Notifier
//
// GoogleDNSNotifier delivers secret expiration alerts to a webhook endpoint
// associated with a Google Cloud project (e.g. a Cloud Function or Cloud Run
// service that manages DNS metadata or records).
//
// # Configuration
//
//	[notifiers.googledns]
//	webhook_url = "https://us-central1-my-project.cloudfunctions.net/vault-alert"
//	project     = "my-gcp-project"
//
// Both webhook_url and project are required.
package notifier
