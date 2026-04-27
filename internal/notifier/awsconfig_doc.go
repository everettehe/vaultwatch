// Package notifier provides notification backends for vaultwatch.
//
// # AWS Config Notifier
//
// The AWSConfigNotifier sends compliance evaluations to AWS Config
// whenever a Vault secret is expiring or has expired.
//
// Secrets that are expired or expiring soon are reported as NON_COMPLIANT;
// all others are reported as COMPLIANT.
//
// # Configuration
//
//	notifiers:
//	  awsconfig:
//	    result_token: "token-from-lambda-event"
//	    region: "us-east-1"
//
// AWS credentials are resolved via the default credential chain
// (environment variables, shared credentials file, IAM role, etc.).
package notifier
