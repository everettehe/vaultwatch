package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func TestBuildNotifiers_Telegram_Valid(t *testing.T) {
	cfg := &config.Config{
		Vault: config.VaultConfig{
			Address: "http://127.0.0.1:8200",
			Token:   "root",
		},
		Notifiers: config.NotifiersConfig{
			Telegram: &config.TelegramConfig{
				BotToken: "bot123token",
				ChatID:   "-100456789",
			},
		},
	}
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_Telegram_MissingChatID(t *testing.T) {
	cfg := &config.Config{
		Vault: config.VaultConfig{
			Address: "http://127.0.0.1:8200",
			Token:   "root",
		},
		Notifiers: config.NotifiersConfig{
			Telegram: &config.TelegramConfig{
				BotToken: "bot123token",
				ChatID:   "",
			},
		},
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing chat ID")
	}
}

func TestTelegramNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewTelegramNotifier("token", "-100123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
