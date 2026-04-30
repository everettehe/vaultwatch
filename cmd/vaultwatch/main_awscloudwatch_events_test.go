package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithCloudWatchEvents() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.CloudWatchEvents = &config.CloudWatchEventsConfig{
		EventBus: "my-event-bus",
		Region:   "us-east-1",
	}
	return cfg
}

func TestBuildNotifiers_CloudWatchEvents_Valid(t *testing.T) {
	cfg := minimalConfigWithCloudWatchEvents()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_CloudWatchEvents_MissingEventBus(t *testing.T) {
	cfg := minimalConfigWithCloudWatchEvents()
	cfg.Notifiers.CloudWatchEvents.EventBus = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing event bus, got nil")
	}
}

func TestBuildNotifiers_CloudWatchEvents_MissingRegion(t *testing.T) {
	cfg := minimalConfigWithCloudWatchEvents()
	cfg.Notifiers.CloudWatchEvents.Region = ""
	// Region is optional (falls back to env/default), so no error expected
	_, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error for missing region, got %v", err)
	}
}

func TestCloudWatchEventsNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewCloudWatchEventsNotifier("bus", "vaultwatch", "SecretExpiration", "us-east-1")
	if err != nil {
		t.Skipf("skipping interface check: %v", err)
	}
	var _ notifier.Notifier = n
}
