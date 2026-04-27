// Package notifier provides notification integrations for VaultWatch.
//
// # SMTP Notifier
//
// The SMTPNotifier sends alert emails via any SMTP server, including the
// AWS SES SMTP interface. It supports PLAIN authentication and TLS.
//
// Configuration example (vaultwatch.yaml):
//
//	notifiers:
//	  smtp:
//	    host: email-smtp.us-east-1.amazonaws.com
//	    port: 587
//	    username: AKIAIOSFODNN7EXAMPLE
//	    password: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
//	    from: alerts@example.com
//	    to:
//	      - oncall@example.com
package notifier
