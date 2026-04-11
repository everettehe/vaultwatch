package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithBearyChat(webhookURL string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.BearyChat.WebhookURL = webhookURL
	return cfg
}

func TestBuildNotifiers_BearyChat_Valid(t *testing.T) {
	cfg := minimalConfigWithBearyChat("https://hook.bearychat.com/abc123")

	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_BearyChat_MissingWebhook(t *testing.T) {
	cfg := minimalConfigWithBearyChat("")

	// No BearyChat webhook configured — should not error, just skip
	_, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error when BearyChat webhook is empty, got %v", err)
	}
}

func TestBearyChatNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewBearyChatNotifier("https://hook.bearychat.com/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var _ notifier.Notifier = n
}
