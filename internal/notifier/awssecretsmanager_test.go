package notifier_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockSMClient struct {
	called bool
	err    error
}

func (m *mockSMClient) PutSecretValue(_ context.Context, _ *secretsmanager.PutSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.PutSecretValueOutput, error) {
	m.called = true
	return &secretsmanager.PutSecretValueOutput{}, m.err
}

func newSMSecret(daysUntil int) vault.Secret {
	return vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour),
	}
}

func TestNewSecretsManagerNotifier_MissingSecretID(t *testing.T) {
	_, err := notifier.NewSecretsManagerNotifier("", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing secret_id")
	}
}

func TestSecretsManagerNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockSMClient{}
	n := notifier.NewSecretsManagerNotifierWithClient(mock, "my-secret", "us-east-1")
	secret := newSMSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected PutSecretValue to be called")
	}
}

func TestSecretsManagerNotifier_Notify_Expired(t *testing.T) {
	mock := &mockSMClient{}
	n := notifier.NewSecretsManagerNotifierWithClient(mock, "my-secret", "us-east-1")
	secret := newSMSecret(-1)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected PutSecretValue to be called")
	}
}

func TestSecretsManagerNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockSMClient{err: errors.New("aws error")}
	n := notifier.NewSecretsManagerNotifierWithClient(mock, "my-secret", "us-east-1")
	secret := newSMSecret(3)
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client")
	}
}

func TestSecretsManagerNotifier_ImplementsInterface(t *testing.T) {
	mock := &mockSMClient{}
	var _ interface {
		Notify(context.Context, vault.Secret) error
	} = notifier.NewSecretsManagerNotifierWithClient(mock, "my-secret", "us-east-1")
}
