// Package notifier provides notification backends for vaultwatch.
//
// # AWS SES Template Notifier
//
// The SESTemplateNotifier sends alerts using pre-configured AWS SES
// email templates. This is useful when you want consistent, branded
// email formatting managed outside of vaultwatch.
//
// Required configuration:
//   - from: sender email address (must be SES-verified)
//   - to: recipient email address (must be SES-verified in sandbox)
//   - template: name of the SES template to use
//   - region: AWS region where the template is registered
//
// Authentication uses the standard AWS credential chain
// (environment variables, shared credentials, IAM role, etc.).
package notifier
