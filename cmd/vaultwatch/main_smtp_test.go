package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithSMTP(host, from string, to []string) *config.Config {
	c := minimalConfig()
	c.Notifiers.SMTP = &config.SMTPConfig{
		Host: host,
		Port: 587,
		From: from,
		To:   to,
	}
	return c
}

func TestBuildNotifiers_SMTP_Valid(t *testing.T) {
	cfg := minimalConfigWithSMTP("smtp.example.com", "from@example.com", []string{"to@example.com"})
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_SMTP_MissingFrom(t *testing.T) {
	cfg := minimalConfigWithSMTP("smtp.example.com", "", []string{"to@example.com"})
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing from address")
	}
}

func TestBuildNotifiers_SMTP_MissingTo(t *testing.T) {
	cfg := minimalConfigWithSMTP("smtp.example.com", "from@example.com", nil)
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing recipients")
	}
}

func TestSMTPNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewSMTPNotifier("smtp.example.com", 587, "", "", "from@example.com", []string{"to@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
