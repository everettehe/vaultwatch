// Package notifier provides alert delivery implementations for vaultwatch.
//
// # VictorOps Notifier
//
// The VictorOpsNotifier sends alerts to VictorOps (now known as Splunk On-Call)
// using the VictorOps REST endpoint integration.
//
// # Configuration
//
// To enable VictorOps alerts, provide the following in your vaultwatch config:
//
//	notifiers:
//	  victorops:
//	    webhook_url: "https://alert.victorops.com/integrations/generic/20131114/alert"
//	    routing_key: "your-routing-key"
//
// # Message Types
//
// Secrets expiring in the future are sent as WARNING messages.
// Already-expired secrets are escalated to CRITICAL.
//
// Each alert includes a stable entity_id derived from the secret path,
// allowing VictorOps to correlate repeated alerts for the same secret.
package notifier
