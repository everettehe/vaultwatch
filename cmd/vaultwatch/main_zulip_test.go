package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithZulip(baseURL, email, apiKey, stream, topic string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.Zulip.BaseURL = baseURL
	cfg.Notifiers.Zulip.BotEmail = email
	cfg.Notifiers.Zulip.APIKey = apiKey
	cfg.Notifiers.Zulip.Stream = stream
	cfg.Notifiers.Zulip.Topic = topic
	return cfg
}

func TestBuildNotifiers_Zulip_Valid(t *testing.T) {
	cfg := minimalConfigWithZulip(
		"https://org.zulipchat.com",
		"bot@org.com",
		"secret-api-key",
		"alerts",
		"vault secrets",
	)
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_Zulip_MissingAPIKey(t *testing.T) {
	cfg := minimalConfigWithZulip(
		"https://org.zulipchat.com",
		"bot@org.com",
		"",
		"alerts",
		"vault",
	)
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestBuildNotifiers_Zulip_MissingStream(t *testing.T) {
	cfg := minimalConfigWithZulip(
		"https://org.zulipchat.com",
		"bot@org.com",
		"apikey",
		"",
		"vault",
	)
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing stream")
	}
}

func TestZulipNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewZulipNotifier(
		"https://org.zulipchat.com",
		"bot@org.com",
		"apikey",
		"alerts",
		"vault",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
