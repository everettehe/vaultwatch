// Package notifier provides alert delivery integrations for VaultWatch.
//
// # Google Cloud Storage Notifier
//
// The GoogleStorageNotifier uploads a JSON alert record to a GCS bucket
// whenever a secret is expiring or has expired. Each notification creates
// a new object whose key encodes the secret path and a UTC timestamp so
// that historical alerts are preserved.
//
// Configuration fields:
//
//	bucket  – GCS bucket name (required)
//	api_key – Google API key with storage.objects.create permission (required)
//
// Example config:
//
//	notifiers:
//	  google_storage:
//	    bucket: "my-vault-alerts"
//	    api_key: "AIza..."
package notifier
