package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func TestBuildNotifiers_Campfire_Valid(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.Campfire.WebhookURL = "https://3.basecamp.com/12345/integrations/abc/buckets/xyz/chats/1/lines"

	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_Campfire_MissingWebhook(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.Campfire.WebhookURL = ""

	// No campfire webhook set — should not error, just skip
	_, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error when campfire webhook is empty, got %v", err)
	}
}

func TestCampfireNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewCampfireNotifier("https://3.basecamp.com/example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}

func minimalConfigWithCampfire() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.Campfire.WebhookURL = "https://3.basecamp.com/12345/integrations/abc/buckets/xyz/chats/1/lines"
	return cfg
}
