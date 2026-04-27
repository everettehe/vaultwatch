// Package notifier provides notification backends for vaultwatch.
//
// SMTPNotifier sends vault secret expiration alerts via SMTP.
// It is compatible with any SMTP server including AWS SES SMTP endpoints.
//
// Configuration example:
//
//	smtp:
//	  host: email-smtp.us-east-1.amazonaws.com
//	  port: 587
//	  username: AKIAIOSFODNN7EXAMPLE
//	  password: your-smtp-password
//	  from: alerts@example.com
//	  to:
//	    - ops@example.com
package notifier
