package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithSESSendTemplate() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.SESSendTemplate = &config.SESSendTemplateConfig{
		From:     "from@example.com",
		To:       "to@example.com",
		Template: "VaultAlert",
		Region:   "us-east-1",
	}
	return cfg
}

func TestBuildNotifiers_SESSendTemplate_MissingFrom(t *testing.T) {
	cfg := minimalConfigWithSESSendTemplate()
	cfg.Notifiers.SESSendTemplate.From = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestBuildNotifiers_SESSendTemplate_MissingTo(t *testing.T) {
	cfg := minimalConfigWithSESSendTemplate()
	cfg.Notifiers.SESSendTemplate.To = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing to")
	}
}

func TestBuildNotifiers_SESSendTemplate_MissingTemplate(t *testing.T) {
	cfg := minimalConfigWithSESSendTemplate()
	cfg.Notifiers.SESSendTemplate.Template = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestSESSendTemplateNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewSESSendTemplateNotifier(
		"from@example.com", "to@example.com", "VaultAlert", "us-east-1",
	)
	if err != nil {
		t.Skipf("skipping interface check (AWS not configured): %v", err)
	}
	var _ notifier.Notifier = n
}
