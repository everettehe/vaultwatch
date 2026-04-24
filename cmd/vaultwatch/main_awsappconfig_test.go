package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithAppConfig() *config.Config {
	cfg := minimalConfig()
	cfg.AppConfig = &config.AppConfigConfig{
		Application: "my-app",
		Environment: "production",
		Profile:     "vault-alerts",
		Region:      "us-east-1",
	}
	return cfg
}

func TestBuildNotifiers_AppConfig_Valid(t *testing.T) {
	cfg := minimalConfigWithAppConfig()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_AppConfig_MissingApplication(t *testing.T) {
	cfg := minimalConfigWithAppConfig()
	cfg.AppConfig.Application = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing application, got nil")
	}
}

func TestBuildNotifiers_AppConfig_MissingEnvironment(t *testing.T) {
	cfg := minimalConfigWithAppConfig()
	cfg.AppConfig.Environment = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing environment, got nil")
	}
}

func TestBuildNotifiers_AppConfig_MissingProfile(t *testing.T) {
	cfg := minimalConfigWithAppConfig()
	cfg.AppConfig.Profile = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}

func TestAppConfigNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewAppConfigNotifier("my-app", "production", "vault-alerts", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
