// Package notifier provides notification backends for VaultWatch.
//
// # AWS Pinpoint Notifier
//
// The Pinpoint notifier sends SMS alerts via AWS Pinpoint when a Vault
// secret is approaching expiration or has already expired.
//
// ## Configuration
//
// The following fields are required in the vaultwatch configuration:
//
//	 pinpoint:
//	   app_id: "abc123def456"        # Pinpoint application/project ID
//	   dest_number: "+15550001234"   # E.164-formatted destination phone number
//	   region: "us-east-1"           # AWS region where the Pinpoint app is deployed
//
// Optional fields:
//
//	   orig_number: "+15559876543"   # Origination number registered in Pinpoint
//	   message_type: "TRANSACTIONAL" # TRANSACTIONAL (default) or PROMOTIONAL
//
// ## IAM Permissions
//
// The IAM principal used by VaultWatch must have the following permission:
//
//	 mobiletargeting:SendMessages
//
// Scope the resource to the specific Pinpoint application ARN where possible:
//
//	 arn:aws:mobiletargeting:<region>:<account-id>:apps/<app-id>
//
// ## Alert Behaviour
//
// A warning SMS is sent when a secret enters the warning threshold window.
// A critical SMS is sent when a secret enters the critical threshold window.
// An expiry SMS is sent once the secret's TTL reaches zero.
//
// Each message includes the secret path and the number of days remaining
// (or a note that the secret has expired).
package notifier
