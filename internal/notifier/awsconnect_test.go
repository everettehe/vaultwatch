package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/connect"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockConnectClient struct {
	called bool
	err    error
}

func (m *mockConnectClient) CreateContact(_ context.Context, _ *connect.CreateContactInput, _ ...func(*connect.Options)) (*connect.CreateContactOutput, error) {
	m.called = true
	return &connect.CreateContactOutput{}, m.err
}

func newConnectSecret(daysUntil int) *vault.Secret {
	expiry := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return &vault.Secret{Path: "secret/db/password", Expiration: expiry}
}

func TestNewAWSConnectNotifier_MissingInstanceID(t *testing.T) {
	_, err := NewAWSConnectNotifier("", "flow-123", "queue-123", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing instance ID")
	}
}

func TestNewAWSConnectNotifier_MissingContactFlow(t *testing.T) {
	_, err := NewAWSConnectNotifier("instance-123", "", "queue-123", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing contact flow")
	}
}

func TestNewAWSConnectNotifier_MissingRegion(t *testing.T) {
	_, err := NewAWSConnectNotifier("instance-123", "flow-123", "queue-123", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewAWSConnectNotifier_Valid(t *testing.T) {
	mock := &mockConnectClient{}
	n, err := newAWSConnectNotifierWithClient(mock, "instance-123", "flow-123", "queue-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestAWSConnectNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockConnectClient{}
	n, _ := newAWSConnectNotifierWithClient(mock, "instance-123", "flow-123", "queue-123")
	secret := newConnectSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Fatal("expected CreateContact to be called")
	}
}

func TestAWSConnectNotifier_Notify_Error(t *testing.T) {
	mock := &mockConnectClient{err: errors.New("connect error")}
	n, _ := newAWSConnectNotifierWithClient(mock, "instance-123", "flow-123", "queue-123")
	secret := newConnectSecret(3)
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from failed CreateContact")
	}
}
