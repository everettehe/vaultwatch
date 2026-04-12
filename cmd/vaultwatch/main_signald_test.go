package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithSignald() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.Signald = &config.SignaldConfig{
		BaseURL:    "http://localhost:8080",
		Sender:     "+15550001111",
		Recipients: []string{"+15550002222"},
	}
	return cfg
}

func TestBuildNotifiers_Signald_Valid(t *testing.T) {
	cfg := minimalConfigWithSignald()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_Signald_MissingSender(t *testing.T) {
	cfg := minimalConfigWithSignald()
	cfg.Notifiers.Signald.Sender = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing sender")
	}
}

func TestBuildNotifiers_Signald_MissingRecipients(t *testing.T) {
	cfg := minimalConfigWithSignald()
	cfg.Notifiers.Signald.Recipients = nil
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing recipients")
	}
}

func TestSignaldNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewSignaldNotifier("http://localhost:8080", "+15550001111", []string{"+15550002222"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
