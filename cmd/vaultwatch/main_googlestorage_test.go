package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithGoogleStorage() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.GoogleStorage = &config.GoogleStorageConfig{
		Bucket: "my-vault-alerts",
		APIKey: "AIza-test-key",
	}
	return cfg
}

func TestBuildNotifiers_GoogleStorage_Valid(t *testing.T) {
	cfg := minimalConfigWithGoogleStorage()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_GoogleStorage_MissingBucket(t *testing.T) {
	cfg := minimalConfigWithGoogleStorage()
	cfg.Notifiers.GoogleStorage.Bucket = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing bucket")
	}
}

func TestBuildNotifiers_GoogleStorage_MissingAPIKey(t *testing.T) {
	cfg := minimalConfigWithGoogleStorage()
	cfg.Notifiers.GoogleStorage.APIKey = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing api_key")
	}
}

func TestGoogleStorageNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewGoogleStorageNotifier("bucket", "key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
