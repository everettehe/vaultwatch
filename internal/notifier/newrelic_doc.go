// Package notifier provides alert delivery integrations for vaultwatch.
//
// # New Relic Notifier
//
// The NewRelicNotifier sends custom events to New Relic Insights using the
// Event API. Each alert is recorded as a "VaultSecretExpiration" event type
// with the following attributes:
//
//   - eventType: always "VaultSecretExpiration"
//   - secretPath: the Vault path of the expiring secret
//   - daysUntilExpiration: floating-point days remaining (negative if expired)
//   - severity: "warning", "critical", or "expired"
//   - message: human-readable description
//   - timestamp: Unix epoch seconds
//
// Configuration fields:
//
//	notifiers:
//	  newrelic:
//	    account_id: "1234567"
//	    api_key: "NRII-..."
//
// The API key must be a New Relic Insert API key (NRII-...).
package notifier
