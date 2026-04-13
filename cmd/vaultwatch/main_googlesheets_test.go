package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func minimalConfigWithGoogleSheets(webAppURL, sheetName string) *config.Config {
	cfg := minimalConfig()
	cfg.GoogleSheets = &config.GoogleSheetsConfig{
		WebAppURL: webAppURL,
		SheetName: sheetName,
	}
	return cfg
}

func TestBuildNotifiers_GoogleSheets_Valid(t *testing.T) {
	cfg := minimalConfigWithGoogleSheets("https://script.google.com/macros/s/abc/exec", "Alerts")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_GoogleSheets_MissingURL(t *testing.T) {
	cfg := minimalConfigWithGoogleSheets("", "")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing web app URL")
	}
}

func TestGoogleSheetsNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewGoogleSheetsNotifier("https://script.google.com/macros/s/abc/exec", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ interface {
		Notify(*vault.Secret) error
	} = n
}
