// Package notifier provides notification backends for vaultwatch.
//
// # SMTP Notifier
//
// The SMTPNotifier sends alert emails via any SMTP server, including
// the AWS SES SMTP interface.
//
// Configuration fields:
//
//	host     - SMTP server hostname (required)
//	port     - SMTP server port (default: 587)
//	username - SMTP auth username (optional)
//	password - SMTP auth password (optional)
//	from     - Sender email address (required)
//	to       - Recipient email address (required)
package notifier
