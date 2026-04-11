package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithCircleCI(token string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.CircleCI.Token = token
	return cfg
}

func TestBuildNotifiers_CircleCI_Valid(t *testing.T) {
	cfg := minimalConfigWithCircleCI("mytoken")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_CircleCI_MissingToken(t *testing.T) {
	cfg := minimalConfigWithCircleCI("")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing CircleCI token")
	}
}

func TestCircleCINotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewCircleCINotifier("tok", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
