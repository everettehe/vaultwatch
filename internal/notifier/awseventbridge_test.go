package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockEventBridgeClient struct {
	called bool
	lastInput *eventbridge.PutEventsInput
	err       error
}

func (m *mockEventBridgeClient) PutEvents(_ context.Context, params *eventbridge.PutEventsInput, _ ...func(*eventbridge.Options)) (*eventbridge.PutEventsOutput, error) {
	m.called = true
	m.lastInput = params
	return &eventbridge.PutEventsOutput{}, m.err
}

func TestNewEventBridgeNotifier_MissingEventBus(t *testing.T) {
	_, err := NewEventBridgeNotifier("", "", "")
	if err == nil {
		t.Fatal("expected error for missing event bus")
	}
}

func TestNewEventBridgeNotifier_Defaults(t *testing.T) {
	mock := &mockEventBridgeClient{}
	n := newEventBridgeNotifierWithClient(mock, "default", "", "")
	if n.source != "vaultwatch" {
		t.Errorf("expected default source 'vaultwatch', got %q", n.source)
	}
	if n.detailType != "VaultSecretExpiry" {
		t.Errorf("expected default detailType 'VaultSecretExpiry', got %q", n.detailType)
	}
}

func TestEventBridgeNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockEventBridgeClient{}
	n := newEventBridgeNotifierWithClient(mock, "default", "vaultwatch", "VaultSecretExpiry")
	s := &vault.Secret{
		Path:      "secret/my-app/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	if err := n.Notify(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Fatal("expected PutEvents to be called")
	}
	if len(mock.lastInput.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(mock.lastInput.Entries))
	}
	entry := mock.lastInput.Entries[0]
	if *entry.Source != "vaultwatch" {
		t.Errorf("unexpected source: %s", *entry.Source)
	}
	if *entry.DetailType != "VaultSecretExpiry" {
		t.Errorf("unexpected detail type: %s", *entry.DetailType)
	}
}

func TestEventBridgeNotifier_Notify_Expired(t *testing.T) {
	mock := &mockEventBridgeClient{}
	n := newEventBridgeNotifierWithClient(mock, "my-bus", "vaultwatch", "VaultSecretExpiry")
	s := &vault.Secret{
		Path:      "secret/expired",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}
	if err := n.Notify(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *mock.lastInput.Entries[0].EventBusName != "my-bus" {
		t.Errorf("unexpected event bus: %s", *mock.lastInput.Entries[0].EventBusName)
	}
}

func TestEventBridgeNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockEventBridgeClient{err: errors.New("aws error")}
	n := newEventBridgeNotifierWithClient(mock, "default", "vaultwatch", "VaultSecretExpiry")
	s := &vault.Secret{
		Path:      "secret/db",
		ExpiresAt: time.Now().Add(72 * time.Hour),
	}
	if err := n.Notify(s); err == nil {
		t.Fatal("expected error from client")
	}
}
