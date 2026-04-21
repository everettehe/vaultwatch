// Package notifier provides alert delivery integrations for VaultWatch.
//
// # AWS EventBridge Notifier
//
// The EventBridge notifier publishes secret expiration events to an Amazon
// EventBridge event bus. Each notification is sent as a custom event with
// a detail-type of "VaultWatch Secret Expiration".
//
// # Configuration
//
//	[notifiers.awseventbridge]
//	event_bus = "my-event-bus"       # required: name or ARN of the event bus
//	source    = "vaultwatch"         # optional: event source (default: "vaultwatch")
//	region    = "us-east-1"          # optional: AWS region
//
// AWS credentials are resolved using the standard AWS credential chain
// (environment variables, shared credentials file, IAM role, etc.).
package notifier
