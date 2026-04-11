package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithLark(webhookURL string) *config.Config {
	cfg := minimalConfig()
	cfg.Lark = &config.LarkConfig{
		WebhookURL: webhookURL,
	}
	return cfg
}

func TestBuildNotifiers_Lark_Valid(t *testing.T) {
	cfg := minimalConfigWithLark("https://open.larksuite.com/open-apis/bot/v2/hook/test")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_Lark_MissingWebhook(t *testing.T) {
	cfg := minimalConfigWithLark("")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing Lark webhook URL")
	}
}

func TestLarkNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewLarkNotifier("https://open.larksuite.com/open-apis/bot/v2/hook/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
