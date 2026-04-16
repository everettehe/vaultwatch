// Package notifier provides alert delivery integrations for VaultWatch.
//
// # AWS SES Notifier
//
// The SES notifier sends email alerts via Amazon Simple Email Service.
//
// Required configuration:
//
//	notifiers:
//	  ses:
//	    from: alerts@example.com
//	    to: ops@example.com
//	    region: us-east-1  # optional, falls back to AWS_REGION env var
//
// Authentication uses the standard AWS credential chain (env vars,
// shared credentials file, IAM role, etc.).
//
// Both the from and to addresses must be verified in SES unless your
// account is out of the SES sandbox.
package notifier
