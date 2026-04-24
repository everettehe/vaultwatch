// Package notifier provides notification implementations for VaultWatch.
//
// # AWS AppConfig Notifier
//
// The AppConfig notifier updates an AWS AppConfig configuration profile
// when a Vault secret is approaching expiration or has expired. This allows
// downstream applications that poll AppConfig to react to secret lifecycle
// events without direct Vault access.
//
// # Configuration
//
//	[notifiers.appconfig]
//	application = "my-app"
//	environment = "production"
//	profile     = "vault-alerts"
//	region      = "us-east-1"
//
// # Required Fields
//
//   - application: the AppConfig application name or ID
//   - environment: the AppConfig environment name or ID
//   - profile:     the AppConfig configuration profile name or ID
//
// # Optional Fields
//
//   - region: AWS region (defaults to AWS_REGION env var or instance metadata)
package notifier
