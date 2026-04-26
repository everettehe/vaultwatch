// Package notifier provides alert delivery integrations for VaultWatch.
//
// # AWS CodeBuild Notifier
//
// The CodeBuild notifier triggers an AWS CodeBuild project build when a
// Vault secret is expiring or has expired. This allows teams to integrate
// secret rotation directly into their CI/CD pipelines.
//
// # Configuration
//
// The following fields are required:
//
//	notifiers:
//	  codebuild:
//	    project_name: "rotate-vault-secrets"
//	    region: "us-east-1"
//
// # IAM Permissions
//
// The AWS credentials used by VaultWatch must have the following permissions:
//
//	codebuild:StartBuild
//
// # Behaviour
//
// When a secret is expiring or expired, a new CodeBuild build is triggered
// for the configured project. The secret path and expiration details are
// passed as environment variable overrides on the build so that the
// CodeBuild project can act on the specific secret.
package notifier
