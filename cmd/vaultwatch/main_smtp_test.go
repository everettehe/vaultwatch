package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithSMTP() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.SMTP = &config.SMTPConfig{
		Host: "smtp.example.com",
		Port: 587,
		From: "from@example.com",
		To:   []string{"to@example.com"},
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
	cfg.Notifiers.SMTP.From = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestBuildNotifiers_SMTP_MissingTo(t *testing.T) {
	cfg := minimalConfigWithSMTP()
	cfg.Notifiers.SMTP.To = nil
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing to")
	}
}

func TestSMTPNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewSMTPNotifier("smtp.example.com", 587, "", "", "from@example.com", []string{"to@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
