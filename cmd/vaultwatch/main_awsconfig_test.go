package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithAWSConfig() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.AWSConfig = &config.AWSConfigConfig{
		ResultToken: "token-abc",
		Region:      "us-east-1",
	}
	return cfg
}

func TestBuildNotifiers_AWSConfig_Valid(t *testing.T) {
	cfg := minimalConfigWithAWSConfig()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_AWSConfig_MissingResultToken(t *testing.T) {
	cfg := minimalConfigWithAWSConfig()
	cfg.Notifiers.AWSConfig.ResultToken = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing result token")
	}
}

func TestBuildNotifiers_AWSConfig_MissingRegion(t *testing.T) {
	cfg := minimalConfigWithAWSConfig()
	cfg.Notifiers.AWSConfig.Region = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestAWSConfigNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewAWSConfigNotifier("token-abc", "us-east-1")
	if err != nil {
		t.Skip("skipping interface check: AWS credentials not available")
	}
	var _ notifier.Notifier = n
}
