package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func TestBuildNotifiers_SignalR_Valid(t *testing.T) {
	cfg := &config.Config{
		SignalR: &config.SignalRConfig{
			EndpointURL: "https://example.service.signalr.net",
			AccessKey:   "mykey",
			Hub:         "alerts",
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

func TestBuildNotifiers_SignalR_MissingAccessKey(t *testing.T) {
	cfg := &config.Config{
		SignalR: &config.SignalRConfig{
			EndpointURL: "https://example.service.signalr.net",
			AccessKey:   "",
			Hub:         "alerts",
		},
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing access key")
	}
}

func TestBuildNotifiers_SignalR_MissingURL(t *testing.T) {
	cfg := &config.Config{
		SignalR: &config.SignalRConfig{
			EndpointURL: "",
			AccessKey:   "mykey",
		},
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing endpoint URL")
	}
}

func TestSignalRNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewSignalRNotifier("https://example.service.signalr.net", "mykey", "alerts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
