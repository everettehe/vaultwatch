package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithGoogleDNS(webhookURL, project string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.GoogleDNS = &config.GoogleDNSConfig{
		WebhookURL: webhookURL,
		Project:    project,
	}
	return cfg
}

func TestBuildNotifiers_GoogleDNS_Valid(t *testing.T) {
	cfg := minimalConfigWithGoogleDNS("https://example.com/hook", "my-project")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_GoogleDNS_MissingProject(t *testing.T) {
	cfg := minimalConfigWithGoogleDNS("https://example.com/hook", "")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestBuildNotifiers_GoogleDNS_MissingWebhook(t *testing.T) {
	cfg := minimalConfigWithGoogleDNS("", "my-project")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestGoogleDNSNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewGoogleDNSNotifier("https://example.com", "proj")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
