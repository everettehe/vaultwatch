package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithGoogleChatCard(webhookURL string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.GoogleChatCard = &config.GoogleChatCardConfig{
		WebhookURL: webhookURL,
	}
	return cfg
}

func TestBuildNotifiers_GoogleChatCard_Valid(t *testing.T) {
	cfg := minimalConfigWithGoogleChatCard("https://chat.googleapis.com/v1/spaces/test")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_GoogleChatCard_MissingWebhook(t *testing.T) {
	cfg := minimalConfigWithGoogleChatCard("")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestGoogleChatCardNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewGoogleChatCardNotifier("https://chat.googleapis.com/v1/spaces/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
