// Package notifier provides alert delivery integrations for vaultwatch.
//
// # Gotify Notifier
//
// The GotifyNotifier sends secret expiration alerts to a self-hosted
// Gotify push notification server (https://gotify.net).
//
// # Configuration
//
//	  gotify:
//	    server_url: "https://gotify.example.com"
//	    token: "your-app-token"
//	    priority: 5  # optional, defaults to 5
//
// # Priority Levels
//
// Gotify uses numeric priority values. Higher values indicate greater urgency.
// A priority of 0 or below will be replaced with the default value of 5.
//
//	  1-3  : low
//	  4-6  : normal (default: 5)
//	  7-9  : high
//	  10+  : urgent
package notifier
