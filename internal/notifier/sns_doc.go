// Package notifier provides alert delivery mechanisms for Vault secret
// expiration events.
//
// SNS Notifier
//
// The SNSNotifier publishes messages to an AWS Simple Notification Service
// (SNS) topic. It relies on the AWS SDK default credential chain, so
// credentials can be supplied via environment variables, shared credential
// files, IAM instance profiles, or any other standard AWS mechanism.
//
// Example configuration in vaultwatch.yaml:
//
//	notifiers:
//	  sns:
//	    topic_arn: "arn:aws:sns:us-east-1:123456789012:vaultwatch-alerts"
//
// Required IAM permission: sns:Publish on the target topic.
package notifier
