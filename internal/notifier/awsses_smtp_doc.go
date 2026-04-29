// Package notifier provides the SMTPNotifier which delivers
// secret-expiration alerts via SMTP. It is compatible with
// AWS SES SMTP credentials as well as any standard SMTP relay.
//
// Configuration fields:
//
//	host     - SMTP server hostname (required)
//	port     - SMTP server port (default: 587)
//	username - SMTP auth username (optional)
//	password - SMTP auth password (optional)
//	from     - sender email address (required)
//	to       - recipient email address (required)
package notifier
