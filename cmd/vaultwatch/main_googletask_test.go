package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithGoogleTask(webhookURL, tasklist string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.GoogleTask = &config.GoogleTaskConfig{
		WebhookURL: webhookURL,
		Tasklist:   tasklist,
	}
	return cfg
}

func TestBuildNotifiers_GoogleTask_Valid(t *testing.T) {
	cfg := minimalConfigWithGoogleTask("https://script.google.com/macros/s/abc/exec", "")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_GoogleTask_MissingWebhook(t *testing.T) {
	cfg := minimalConfigWithGoogleTask("", "")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestGoogleTaskNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewGoogleTaskNotifier("https://example.com", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
