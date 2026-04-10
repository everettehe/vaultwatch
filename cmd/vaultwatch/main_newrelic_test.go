package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func TestBuildNotifiers_NewRelic_Valid(t *testing.T) {
	cfg := &config.Config{
		NewRelic: &config.NewRelicConfig{
			AccountID: "123456",
			APIKey:    "NRII-test-key",
		},
	}
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_NewRelic_MissingAPIKey(t *testing.T) {
	cfg := &config.Config{
		NewRelic: &config.NewRelicConfig{
			AccountID: "123456",
			APIKey:    "",
		},
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestBuildNotifiers_NewRelic_MissingAccountID(t *testing.T) {
	cfg := &config.Config{
		NewRelic: &config.NewRelicConfig{
			AccountID: "",
			APIKey:    "NRII-test-key",
		},
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing account ID")
	}
}

func TestNewRelicNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewNewRelicNotifier("123456", "NRII-test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
