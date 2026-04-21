package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithEventBridge(eventBus string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.AWSEventBridge = &config.AWSEventBridgeConfig{
		EventBus: eventBus,
		Source:   "vaultwatch",
	}
	return cfg
}

func TestBuildNotifiers_EventBridge_Valid(t *testing.T) {
	cfg := minimalConfigWithEventBridge("my-event-bus")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_EventBridge_MissingEventBus(t *testing.T) {
	cfg := minimalConfigWithEventBridge("")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing event bus, got nil")
	}
}

func TestEventBridgeNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewEventBridgeNotifier("test-bus", "vaultwatch", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
