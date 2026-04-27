package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithSMTP() *config.Config {
	cfg := minimalConfig()
	cfg.SMTP = &config.SMTPConfig{
		Host: "smtp.example.com",
		From: "from@example.com",
		To:   "to@example.com",
	}
	return cfg
}

func TestBuildNotifiers_SMTP_Valid(t *testing.T) {
	cfg := minimalConfigWithSMTP()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_SMTP_MissingFrom(t *testing.T) {
	cfg := minimalConfigWithSMTP()
	cfg.SMTP.From = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing from address")
	}
}

func TestBuildNotifiers_SMTP_MissingTo(t *testing.T) {
	cfg := minimalConfigWithSMTP()
	cfg.SMTP.To = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing to address")
	}
}

func TestSMTPNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewSMTPNotifier("smtp.example.com", "", "", "from@example.com", "to@example.com", 587)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
