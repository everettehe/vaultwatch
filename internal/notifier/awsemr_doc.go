// Package notifier provides alert delivery integrations for vaultwatch.
//
// # AWS EMR Notifier
//
// The EMRNotifier submits a lightweight EMR job flow step to an existing
// Amazon EMR cluster whenever a Vault secret is approaching expiration or
// has already expired. The step uses the built-in "command-runner.jar" to
// echo a JSON payload containing the secret path, expiration timestamp, and
// a human-readable message.
//
// # Configuration
//
//	 notifiers:
//	   emr:
//	     cluster_id: "j-XXXXXXXXXXXXX"
//	     region: "us-east-1"
//
// AWS credentials are resolved via the default credential chain
// (environment variables, ~/.aws/credentials, IAM role, etc.).
package notifier
