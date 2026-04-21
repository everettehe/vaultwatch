package main

import (
	"context"
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func minimalConfigWithBigQuery() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.BigQuery = &config.BigQueryConfig{
		ProjectID: "my-project",
		DatasetID: "vaultwatch",
		TableID:   "expirations",
		APIKey:    "test-api-key",
	}
	return cfg
}

func TestBuildNotifiers_BigQuery_Valid(t *testing.T) {
	cfg := minimalConfigWithBigQuery()
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_BigQuery_MissingAPIKey(t *testing.T) {
	cfg := minimalConfigWithBigQuery()
	cfg.Notifiers.BigQuery.APIKey = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing api_key")
	}
}

func TestBuildNotifiers_BigQuery_MissingProjectID(t *testing.T) {
	cfg := minimalConfigWithBigQuery()
	cfg.Notifiers.BigQuery.ProjectID = ""
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing project_id")
	}
}

func TestBigQueryNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewBigQueryNotifier("proj", "ds", "tbl", "key")
	if err != nil {
		t.Fatal(err)
	}
	var _ interface {
		Notify(context.Context, vault.Secret) error
	} = n
}
