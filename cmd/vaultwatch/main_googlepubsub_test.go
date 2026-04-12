package main

import (
	"testing"

	"github.com/warden-protocol/vaultwatch/internal/config"
	"github.com/warden-protocol/vaultwatch/internal/notifier"
)

func minimalConfigWithGooglePubSub(projectID, topicID string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.GooglePubSub = &config.GooglePubSubConfig{
		ProjectID: projectID,
		TopicID:   topicID,
	}
	return cfg
}

func TestBuildNotifiers_GooglePubSub_MissingProjectID(t *testing.T) {
	cfg := minimalConfigWithGooglePubSub("", "vault-alerts")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing project_id, got nil")
	}
}

func TestBuildNotifiers_GooglePubSub_MissingTopicID(t *testing.T) {
	cfg := minimalConfigWithGooglePubSub("my-project", "")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing topic_id, got nil")
	}
}

func TestGooglePubSubNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewGooglePubSubNotifier("my-project", "vault-alerts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
