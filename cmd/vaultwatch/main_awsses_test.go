package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithSES() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.SES = &config.SESConfig{
		From:   "alerts@example.com",
		To:     "ops@example.com",
		Region: "us-east-1",
	}
	return cfg
}

func TestBuildNotifiers_SES_MissingFrom(t *testing.T) {
	cfg := minimalConfigWithSES()
	cfg.Notifiers.SES.From = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing SES from address")
	}
}

func TestBuildNotifiers_SES_MissingTo(t *testing.T) {
	cfg := minimalConfigWithSES()
	cfg.Notifiers.SES.To = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing SES to address")
	}
}

func TestSESNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewSESNotifier("alerts@example.com", "ops@example.com", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
