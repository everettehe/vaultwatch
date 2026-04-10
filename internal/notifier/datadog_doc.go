// Package notifier provides alert delivery implementations for vaultwatch.
//
// # Datadog Notifier
//
// The DatadogNotifier sends secret expiration events to the Datadog Events API
// (https://docs.datadoghq.com/api/latest/events/).
//
// Each event includes:
//   - A formatted subject and body via FormatMessage
//   - An alert_type of "warning" for expiring secrets or "error" for expired ones
//   - Tags: source:vaultwatch and secret_path:<path>
//
// Configuration requires a Datadog API key. The API URL defaults to the
// standard Datadog US endpoint but can be overridden for EU regions or testing.
//
// Example usage:
//
//	notifier, err := notifier.NewDatadogNotifier(os.Getenv("DD_API_KEY"), "")
//	if err != nil {
//		log.Fatal(err)
//	}
package notifier
