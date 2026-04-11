package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func TestBuildNotifiers_Gotify_Valid(t *testing.T) {
	cfg := &config.Config{
		Gotify: &config.GotifyConfig{
			ServerURL: "http://gotify.example.com",
			Token:     "apptoken",
			Priority:  5,
		},
	}
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_Gotify_MissingToken(t *testing.T) {
	cfg := &config.Config{
		Gotify: &config.GotifyConfig{
			ServerURL: "http://gotify.example.com",
			Token:     "",
		},
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestBuildNotifiers_Gotify_MissingURL(t *testing.T) {
	cfg := &config.Config{
		Gotify: &config.GotifyConfig{
			ServerURL: "",
			Token:     "apptoken",
		},
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing server URL")
	}
}

func TestGotifyNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewGotifyNotifier("http://gotify.example.com", "tok", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
