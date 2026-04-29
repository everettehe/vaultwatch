package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithSESRaw() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.SESRaw = &config.SESRawConfig{
		From:   "from@example.com",
		To:     "to@example.com",
		Region: "us-east-1",
	}
	return cfg
}

func TestBuildNotifiers_SESRaw_MissingFrom(t *testing.T) {
	cfg := minimalConfigWithSESRaw()
	cfg.Notifiers.SESRaw.From = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing from address")
	}
}

func TestBuildNotifiers_SESRaw_MissingTo(t *testing.T) {
	cfg := minimalConfigWithSESRaw()
	cfg.Notifiers.SESRaw.To = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing to address")
	}
}

func TestSESRawNotifier_ImplementsInterface(t *testing.T) {
	n := notifier.newSESRawNotifierWithClient(nil, "from@example.com", "to@example.com")
	var _ notifier.Notifier = n
}
