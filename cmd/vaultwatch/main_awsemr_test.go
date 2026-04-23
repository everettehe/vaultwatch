package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithEMR(clusterID, region string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.EMR = &config.EMRConfig{
		ClusterID: clusterID,
		Region:    region,
	}
	return cfg
}

func TestBuildNotifiers_EMR_Valid(t *testing.T) {
	cfg := minimalConfigWithEMR("j-XXXXXXXXXXXXX", "us-east-1")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_EMR_MissingClusterID(t *testing.T) {
	cfg := minimalConfigWithEMR("", "us-east-1")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing cluster ID")
	}
}

func TestBuildNotifiers_EMR_MissingRegion(t *testing.T) {
	cfg := minimalConfigWithEMR("j-XXXXXXXXXXXXX", "")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestEMRNotifier_ImplementsInterface(t *testing.T) {
	client := &struct{ notifier.EMRNotifier }{}
	_ = client
	var _ notifier.Notifier = (*notifier.EMRNotifier)(nil)
}
