package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithGoogleChat(webhookURL string) *config.Config {
	cfg := minimalConfig()
	cfg.GoogleChat = &config.GoogleChatConfig{
		WebhookURL: webhookURL,
	}
	return cfg
}

func TestBuildNotifiers_GoogleChat_Valid(t *testing.T) {
	cfg := minimalConfigWithGoogleChat("https://chat.googleapis.com/v1/spaces/xxx/messages?key=yyy")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_GoogleChat_MissingWebhook(t *testing.T) {
	cfg := minimalConfigWithGoogleChat("")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestGoogleChatNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewGoogleChatNotifier("https://chat.googleapis.com/v1/spaces/xxx/messages?key=yyy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
