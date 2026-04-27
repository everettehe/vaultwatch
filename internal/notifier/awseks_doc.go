// Package notifier provides notification backends for VaultWatch.
//
// # AWS EKS Notifier
//
// The EKS notifier triggers an EKS managed node group update or addon update
// when a Vault secret is expiring or has expired. This can be used to
// automatically roll workloads that depend on the secret.
//
// Configuration example:
//
//	notifiers:
//	  eks:
//	    cluster: my-cluster
//	    addon: coredns
//	    region: us-east-1
//
// The notifier requires standard AWS credentials to be available via
// environment variables, EC2 instance profile, or shared credentials file.
package notifier
