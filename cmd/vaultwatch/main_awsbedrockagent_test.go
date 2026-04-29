package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithBedrockAgent() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.BedrockAgent = &config.BedrockAgentConfig{
		AgentID: "agent-abc",
		AliasID: "alias-xyz",
		Region:  "us-east-1",
	}
	return cfg
}

func TestBuildNotifiers_BedrockAgent_Valid(t *testing.T) {
	cfg := minimalConfigWithBedrockAgent()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_BedrockAgent_MissingAgentID(t *testing.T) {
	cfg := minimalConfigWithBedrockAgent()
	cfg.Notifiers.BedrockAgent.AgentID = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing agent_id")
	}
}

func TestBuildNotifiers_BedrockAgent_MissingRegion(t *testing.T) {
	cfg := minimalConfigWithBedrockAgent()
	cfg.Notifiers.BedrockAgent.Region = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestBedrockAgentNotifier_ImplementsInterface(t *testing.T) {
	var _ notifier.Notifier = (*notifier.BedrockAgentNotifier)(nil)
}
