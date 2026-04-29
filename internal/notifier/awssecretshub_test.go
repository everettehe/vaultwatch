package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockSecurityHubClient struct {
	called bool
	err    error
}

func (m *mockSecurityHubClient) BatchImportFindings(_ context.Context, _ *securityhub.BatchImportFindingsInput, _ ...func(*securityhub.Options)) (*securityhub.BatchImportFindingsOutput, error) {
	m.called = true
	return &securityhub.BatchImportFindingsOutput{}, m.err
}

func newSecurityHubSecret(daysUntil int) *vault.Secret {
	expiry := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return &vault.Secret{Path: "secret/myapp/api-key", Expiration: expiry}
}

func TestNewSecurityHubNotifier_MissingAccountID(t *testing.T) {
	_, err := NewSecurityHubNotifier("", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing account ID")
	}
}

func TestNewSecurityHubNotifier_MissingRegion(t *testing.T) {
	_, err := NewSecurityHubNotifier("123456789012", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewSecurityHubNotifier_Valid(t *testing.T) {
	client := &mockSecurityHubClient{}
	n := newSecurityHubNotifierWithClient(client, "123456789012", "us-east-1")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	if n.accountID != "123456789012" {
		t.Errorf("expected accountID 123456789012, got %s", n.accountID)
	}
}

func TestSecurityHubNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockSecurityHubClient{}
	n := newSecurityHubNotifierWithClient(client, "123456789012", "us-east-1")
	s := newSecurityHubSecret(5)

	if err := n.Notify(context.Background(), s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Error("expected BatchImportFindings to be called")
	}
}

func TestSecurityHubNotifier_Notify_Expired(t *testing.T) {
	client := &mockSecurityHubClient{}
	n := newSecurityHubNotifierWithClient(client, "123456789012", "us-east-1")
	s := newSecurityHubSecret(-1)

	if err := n.Notify(context.Background(), s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Error("expected BatchImportFindings to be called")
	}
}

func TestSecurityHubNotifier_Notify_ClientError(t *testing.T) {
	client := &mockSecurityHubClient{err: errors.New("aws error")}
	n := newSecurityHubNotifierWithClient(client, "123456789012", "us-east-1")
	s := newSecurityHubSecret(10)

	if err := n.Notify(context.Background(), s); err == nil {
		t.Fatal("expected error from client")
	}
}
