package main

import (
	"os"
	"testing"
)

func TestRun_MissingConfig(t *testing.T) {
	t.Setenv("VAULTWATCH_CONFIG", "/nonexistent/path/config.yaml")

	err := run()
	if err == nil {
		t.Fatal("expected error for missing config file, got nil")
	}
}

func TestRun_InvalidConfig(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "vaultwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer f.Close()

	_, _ = f.WriteString("vault:\n  address: \"\"\n")
	t.Setenv("VAULTWATCH_CONFIG", f.Name())

	err = run()
	if err == nil {
		t.Fatal("expected error for invalid config, got nil")
	}
}

func TestBuildNotifiers_LogOnly(t *testing.T) {
	cfg := minimalConfig()

	n, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestBuildNotifiers_InvalidSlack(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.Slack.WebhookURL = "not-a-url"

	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for invalid slack webhook, got nil")
	}
}
