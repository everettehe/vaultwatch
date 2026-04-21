// Package notifier provides alert delivery integrations for vaultwatch.
//
// # Google Task Notifier
//
// The GoogleTaskNotifier creates tasks in Google Tasks when a Vault secret
// is approaching expiration or has already expired. It posts a JSON payload
// to a configured webhook URL (e.g., a Google Apps Script web app that
// wraps the Tasks API).
//
// # Configuration
//
//	 notifiers:
//	   googletask:
//	     webhook_url: "https://script.google.com/macros/s/<ID>/exec"
//	     tasklist: "@default"   # optional, defaults to @default
//
// # Payload
//
// The notifier sends a JSON body with the following fields:
//
//	{
//	  "title": "<alert subject>",
//	  "notes": "<alert body>",
//	  "due":   "<RFC3339 timestamp>",
//	  "status": "needsAction"
//	}
package notifier
