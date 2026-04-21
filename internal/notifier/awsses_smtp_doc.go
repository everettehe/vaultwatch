// Package notifier provides alert delivery integrations for VaultWatch.
//
// # SMTP Notifier
//
// The SMTP notifier sends email alerts via a plain SMTP server (e.g. AWS SES
// SMTP interface, SendGrid, Mailgun, or any self-hosted MTA).
//
// Configuration example:
//
//	smtp:
//	  host: email-smtp.us-east-1.amazonaws.com
//	  port: 587
//	  username: AKIAIOSFODNN7EXAMPLE
//	  password: secret
//	  from: vaultwatch@example.com
//	  to: ops@example.com
//	  tls: true
//
// The notifier uses STARTTLS when tls is true (default: true).
package notifier
