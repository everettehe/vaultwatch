//go:build integration
// +build integration

package notifier_test

import (
	"context"
	"testing"
	"time"

	"github.com/warden-protocol/vaultwatch/internal/notifier"
	"github.com/warden-protocol/vaultwatch/internal/vault"
)

// TestGooglePubSubNotifier_Integration publishes a real message to Pub/Sub.
// Requires GOOGLE_APPLICATION_CREDENTIALS and a running emulator or real topic.
// Run with: go test -tags integration ./internal/notifier/...
func TestGooglePubSubNotifier_Integration(t *testing.T) {
	projectID := "my-gcp-project"
	topicID := "vault-alerts-test"

	n, err := notifier.NewGooglePubSubNotifier(projectID, topicID)
	if err != nil {
		t.Fatalf("failed to create notifier: %v", err)
	}

	s := &vault.Secret{
		Path:      "secret/data/db/password",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	ctx := context.Background()
	if err := n.Notify(ctx, s); err != nil {
		t.Fatalf("Notify failed: %v", err)
	}
}
