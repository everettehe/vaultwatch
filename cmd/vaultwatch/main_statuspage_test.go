package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithStatusPage() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.StatusPage = &config.StatusPageConfig{
		APIKey: "key123",
		PageID: "page456",
	}
	return cfg
}

func TestBuildNotifiers_StatusPage_Valid(t *testing.T) {
	cfg := minimalConfigWithStatusPage()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_StatusPage_MissingAPIKey(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.StatusPage = &config.StatusPageConfig{
		APIKey: "",
		PageID: "page456",
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing api key")
	}
}

func TestBuildNotifiers_StatusPage_MissingPageID(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.StatusPage = &config.StatusPageConfig{
		APIKey: "key123",
		PageID: "",
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing page ID")
	}
}

func TestStatusPageNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewStatusPageNotifier("key", "page")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
