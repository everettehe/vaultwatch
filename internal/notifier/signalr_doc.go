// Package notifier provides alert delivery implementations for vaultwatch.
//
// # SignalR Notifier
//
// The SignalRNotifier delivers secret expiration alerts to an Azure SignalR
// Service endpoint. It posts a JSON payload to the configured hub using
// Bearer token authentication.
//
// # Configuration
//
//	[signalr]
//	endpoint_url = "https://<resource>.service.signalr.net"
//	access_key   = "<your-access-key>"
//	hub          = "vaultwatch"   # optional, defaults to "vaultwatch"
//
// # Payload Format
//
// The notifier posts to: POST {endpoint_url}/api/v1/hubs/{hub}
//
//	{
//	  "target": "vaultAlert",
//	  "arguments": ["<subject>", "<body>"]
//	}
package notifier
