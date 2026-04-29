// Package notifier provides notification integrations for vaultwatch.
//
// SMTPNotifier sends secret expiration alerts via SMTP. It is compatible
// with AWS SES SMTP endpoints as well as any standard SMTP server.
//
// Configuration example (vaultwatch.yaml):
//
//	notifiers:
//	  smtp:
//	    host: email-smtp.us-east-1.amazonaws.com
//	    port: 587
//	    username: AKIAIOSFODNN7EXAMPLE
//	    password: your-smtp-password
//	    from: alerts@example.com
//	    to:
//	      - oncall@example.com
package notifier
