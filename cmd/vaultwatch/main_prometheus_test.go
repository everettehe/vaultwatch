package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func TestBuildNotifiers_Prometheus_Valid(t *testing.T) {
	cfg := &config.Config{
		Prometheus: &config.PrometheusConfig{
			PushgatewayURL: "http://localhost:9091",
			Job:            "vaultwatch",
		},
	}

	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_Prometheus_MissingURL(t *testing.T) {
	cfg := &config.Config{
		Prometheus: &config.PrometheusConfig{
			PushgatewayURL: "",
			Job:            "vaultwatch",
		},
	}

	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing pushgateway URL")
	}
}

func TestPrometheusNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewPrometheusNotifier("http://localhost:9091", "vaultwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ notifier.Notifier = n
}
