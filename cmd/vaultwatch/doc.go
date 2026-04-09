// Package main is the entry point for the vaultwatch CLI tool.
//
// vaultwatch monitors HashiCorp Vault secret expiration and sends
// configurable alerts before rotation deadlines. It reads configuration
// from a YAML file (default: configs/vaultwatch.yaml) or from the path
// specified by the VAULTWATCH_CONFIG environment variable.
//
// Usage:
//
//	# Run with default config path
//	vaultwatch
//
//	# Run with custom config path
//	VAULTWATCH_CONFIG=/etc/vaultwatch/config.yaml vaultwatch
//
// The tool will scan configured Vault paths for secrets, check their
// expiration metadata, and dispatch notifications via the configured
// channels (log, Slack, email) when secrets are expiring soon or expired.
package main
