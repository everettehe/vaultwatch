package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func TestBuildNotifiers_VictorOps_Valid(t *testing.T) {
	cfg := &config.Config{
		Vault: config.VaultConfig{
			Address: "http://127.0.0.1:8200",
			Token:   "root",
		},
		Notifiers: config.NotifiersConfig{
			VictorOps: &config.VictorOpsConfig{
				WebhookURL: "https://alert.victorops.com/integrations/generic",
				RoutingKey: "my-key",
			},
		},
	}

	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// log notifier + victorops notifier
	if len(notifiers) < 2 {
		t.Fatalf("expected at least 2 notifiers, got %d", len(notifiers))
	}
}

func TestBuildNotifiers_VictorOps_MissingRoutingKey(t *testing.T) {
	cfg := &config.Config{
		Vault: config.VaultConfig{
			Address: "http://127.0.0.1:8200",
			Token:   "root",
		},
		Notifiers: config.NotifiersConfig{
			VictorOps: &config.VictorOpsConfig{
				WebhookURL: "https://alert.victorops.com/integrations/generic",
				RoutingKey: "",
			},
		},
	}

	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing routing key")
	}
}

func TestVictorOpsNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewVictorOpsNotifier("https://example.com", "key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
