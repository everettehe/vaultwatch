package main

import (
	"context"
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func minimalConfigWithSecretsManager() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.AWSSecretsManager = &config.AWSSecretsManagerConfig{
		SecretID: "arn:aws:secretsmanager:us-east-1:123456789012:secret:vaultwatch",
		Region:   "us-east-1",
	}
	return cfg
}

func TestBuildNotifiers_SecretsManager_MissingSecretID(t *testing.T) {
	cfg := minimalConfig()
	cfg.Notifiers.AWSSecretsManager = &config.AWSSecretsManagerConfig{
		Region: "us-east-1",
	}
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error when secret_id is missing")
	}
}

func TestSecretsManagerNotifier_ImplementsInterface(t *testing.T) {
	type notifierIface interface {
		Notify(context.Context, vault.Secret) error
	}
	var _ notifierIface = (*notifier.SecretsManagerNotifier)(nil)
}
