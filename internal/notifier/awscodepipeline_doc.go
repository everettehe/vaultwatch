// Package notifier provides notification backends for vaultwatch.
//
// # AWS CodePipeline Notifier
//
// The CodePipelineNotifier sends a PutJobFailureResult call to AWS CodePipeline
// when a monitored Vault secret is expiring or has expired. This is useful when
// Vault secret rotation is gated behind a pipeline stage and you want the pipeline
// to reflect the unhealthy state of a credential.
//
// # Configuration
//
//	notifiers:
//	  codepipeline:
//	    job_id: "<pipeline-job-id>"
//	    region: "us-east-1"
//
// AWS credentials are resolved via the standard AWS credential chain
// (environment variables, shared credentials file, IAM role, etc.).
package notifier
