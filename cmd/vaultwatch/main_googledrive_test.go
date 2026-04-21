package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithGoogleDrive() *config.Config {
	cfg := minimalConfig()
	cfg.GoogleDrive = &config.GoogleDriveConfig{
		SpreadsheetID: "sheet123",
		SheetName:     "Alerts",
		APIKey:        "test-api-key",
	}
	return cfg
}

func TestBuildNotifiers_GoogleDrive_Valid(t *testing.T) {
	cfg := minimalConfigWithGoogleDrive()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_GoogleDrive_MissingAPIKey(t *testing.T) {
	cfg := minimalConfigWithGoogleDrive()
	cfg.GoogleDrive.APIKey = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error when API key is missing")
	}
}

func TestBuildNotifiers_GoogleDrive_MissingSpreadsheetID(t *testing.T) {
	cfg := minimalConfigWithGoogleDrive()
	cfg.GoogleDrive.SpreadsheetID = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error when spreadsheet ID is missing")
	}
}

func TestGoogleDriveNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewGoogleDriveNotifier("sheet123", "Alerts", "key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
