package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithClickSend(username, apiKey, to string) *config.Config {
	cfg := minimalConfig()
	cfg.ClickSend = &config.ClickSendConfig{
		Username: username,
		APIKey:   apiKey,
		To:       to,
	}
	return cfg
}

func TestBuildNotifiers_ClickSend_Valid(t *testing.T) {
	cfg := minimalConfigWithClickSend("user", "key123", "+15551234567")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_ClickSend_MissingAPIKey(t *testing.T) {
	cfg := minimalConfigWithClickSend("user", "", "+15551234567")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing api key")
	}
}

func TestBuildNotifiers_ClickSend_MissingTo(t *testing.T) {
	cfg := minimalConfigWithClickSend("user", "key123", "")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing recipient")
	}
}

func TestClickSendNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewClickSendNotifier("user", "key123", "+15551234567")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
