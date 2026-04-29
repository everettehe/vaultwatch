package config

// BedrockAgentConfig holds configuration for the AWS Bedrock Agent notifier.
type BedrockAgentConfig struct {
	AgentID  string `mapstructure:"agent_id"`
	AliasID  string `mapstructure:"alias_id"`
	Region   string `mapstructure:"region"`
}
