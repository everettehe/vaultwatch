package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func TestBuildNotifiers_Pushover_Valid(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.Pushover.UserKey = "userkey123"
	cfg.Notifiers.Pushover.APIToken = "apitoken456"

	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_Pushover_MissingAPIToken(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.Pushover.UserKey = "userkey123"
	cfg.Notifiers.Pushover.APIToken = ""

	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing Pushover API token")
	}
}

func TestBuildNotifiers_Pushover_MissingUserKey(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.Pushover.UserKey = ""
	cfg.Notifiers.Pushover.APIToken = "apitoken456"

	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing Pushover user key")
	}
}

func TestBuildNotifiers_Pushover_MissingBothCredentials(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.Pushover.UserKey = ""
	cfg.Notifiers.Pushover.APIToken = ""

	// When both credentials are missing, no Pushover notifier should be
	// registered and no error should be returned (notifier simply skipped).
	_, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error when both Pushover credentials are absent, got %v", err)
	}
}

func TestPushoverNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewPushoverNotifier("userkey123", "apitoken456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}

func minimalConfigWithPushover() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.Pushover.UserKey = "userkey123"
	cfg.Notifiers.Pushover.APIToken = "apitoken456"
	return cfg
}
