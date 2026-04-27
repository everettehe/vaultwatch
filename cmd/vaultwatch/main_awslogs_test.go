package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithCloudWatchLogsV2() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.CloudWatchLogs = &config.CloudWatchLogsConfig{
		LogGroup: "vaultwatch-logs",
		Region:   "us-east-1",
	}
	return cfg
}

func TestBuildNotifiers_CloudWatchLogsV2_Valid(t *testing.T) {
	cfg := minimalConfigWithCloudWatchLogsV2()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_CloudWatchLogsV2_MissingLogGroup(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.CloudWatchLogs = &config.CloudWatchLogsConfig{
		Region: "us-east-1",
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing log group")
	}
}

func TestCloudWatchLogsV2Notifier_ImplementsInterface(t *testing.T) {
	cfg := config.CloudWatchLogsConfig{
		LogGroup: "vaultwatch-logs",
		Region:   "us-east-1",
	}
	n, err := notifier.NewCloudWatchLogsV2Notifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
