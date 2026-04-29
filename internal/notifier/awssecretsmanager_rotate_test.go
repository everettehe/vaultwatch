package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	vaultsecret "github.com/youorg/vaultwatch/internal/vault"
)

type mockRotateClient struct {
	called bool
	err    error
}

func (m *mockRotateClient) RotateSecret(_ context.Context, _ *secretsmanager.RotateSecretInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.RotateSecretOutput, error) {
	m.called = true
	return &secretsmanager.RotateSecretOutput{}, m.err
}

func newRotateSecret() *vaultsecret.Secret {
	return &vaultsecret.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
}

func TestNewSecretsManagerRotateNotifier_MissingSecretID(t *testing.T) {
	_, err := NewSecretsManagerRotateNotifier("", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing secret_id")
	}
}

func TestNewSecretsManagerRotateNotifier_MissingRegion(t *testing.T) {
	_, err := NewSecretsManagerRotateNotifier("arn:aws:secretsmanager:us-east-1:123456789012:secret:test", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestSecretsManagerRotateNotifier_Notify_Success(t *testing.T) {
	mc := &mockRotateClient{}
	n := newSecretsManagerRotateNotifierWithClient("my-secret", mc)

	if err := n.Notify(context.Background(), newRotateSecret()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mc.called {
		t.Error("expected RotateSecret to be called")
	}
}

func TestSecretsManagerRotateNotifier_Notify_Error(t *testing.T) {
	mc := &mockRotateClient{err: errors.New("rotation failed")}
	n := newSecretsManagerRotateNotifierWithClient("my-secret", mc)

	err := n.Notify(context.Background(), newRotateSecret())
	if err == nil {
		t.Fatal("expected error from RotateSecret")
	}
}

func TestSecretsManagerRotateNotifier_ImplementsInterface(t *testing.T) {
	mc := &mockRotateClient{}
	n := newSecretsManagerRotateNotifierWithClient("my-secret", mc)
	var _ Notifier = n
}
