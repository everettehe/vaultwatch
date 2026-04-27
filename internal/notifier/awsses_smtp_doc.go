// Package notifier provides notification backends for vaultwatch.
//
// # SMTP Notifier
//
// The SMTPNotifier sends alert emails directly via an SMTP server.
// It is compatible with AWS SES SMTP endpoints as well as any
// standard SMTP relay.
//
// # Configuration
//
//	smtp:
//	  host: email-smtp.us-east-1.amazonaws.com
//	  port: 587
//	  username: AKIAIOSFODNN7EXAMPLE
//	  password: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
//	  from: alerts@example.com
//	  to: oncall@example.com
package notifier
