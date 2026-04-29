// Package notifier provides notification backends for vaultwatch.
//
// # AWS Bedrock Agent Notifier
//
// The BedrockAgentNotifier sends secret expiry events as natural-language
// prompts to an AWS Bedrock Agent. This allows you to wire expiry alerts
// into AI-driven remediation workflows.
//
// # Configuration
//
//	[notifiers.bedrock_agent]
//	agent_id = "ABCDE12345"
//	alias_id = "TSTALIASID"
//	region   = "us-east-1"
//
// # Required Fields
//
//   - agent_id:  The Bedrock Agent resource ID.
//   - alias_id:  The Bedrock Agent alias ID.
//   - region:    AWS region where the agent is deployed.
package notifier
