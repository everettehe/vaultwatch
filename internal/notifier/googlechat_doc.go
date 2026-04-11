// Package notifier provides alert delivery mechanisms for VaultWatch.
//
// # Google Chat Notifier
//
// The Google Chat notifier sends secret expiration alerts to a Google Chat
// space via an incoming webhook URL.
//
// # Configuration
//
// Set the following fields in your vaultwatch.yaml:
//
//	notifiers:
//	  googlechat:
//	    webhook_url: "https://chat.googleapis.com/v1/spaces/.../messages?key=...&token=..."
//
// # Behavior
//
// Each alert is sent as a simple text message containing the secret path,
// days until expiration, and severity level. Both warning and critical
// thresholds trigger a notification.
package notifier
