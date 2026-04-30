package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockWAFClient struct {
	called bool
	err    error
}

func (m *mockWAFClient) UpdateIPSet(_ context.Context, _ *wafv2.UpdateIPSetInput, _ ...func(*wafv2.Options)) (*wafv2.UpdateIPSetOutput, error) {
	m.called = true
	return &wafv2.UpdateIPSetOutput{}, m.err
}

func newWAFSecret(daysLeft int) *vault.Secret {
	expiry := time.Now().Add(time.Duration(daysLeft) * 24 * time.Hour)
	return &vault.Secret{Path: "secret/waf-test", ExpiresAt: expiry}
}

func TestNewWAFNotifier_MissingIPSetID(t *testing.T) {
	_, err := NewWAFNotifier("", "my-set", "REGIONAL", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing ip_set_id")
	}
}

func TestNewWAFNotifier_MissingRegion(t *testing.T) {
	_, err := NewWAFNotifier("abc123", "my-set", "REGIONAL", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewWAFNotifier_Valid(t *testing.T) {
	client := &mockWAFClient{}
	n := newWAFNotifierWithClient(client, "set-id", "set-name", types.ScopeRegional, "us-east-1")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestWAFNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockWAFClient{}
	n := newWAFNotifierWithClient(client, "set-id", "set-name", types.ScopeRegional, "us-east-1")
	secret := newWAFSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Error("expected UpdateIPSet to be called")
	}
}

func TestWAFNotifier_Notify_Expired(t *testing.T) {
	client := &mockWAFClient{}
	n := newWAFNotifierWithClient(client, "set-id", "set-name", types.ScopeCloudfront, "us-east-1")
	secret := newWAFSecret(-2)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Error("expected UpdateIPSet to be called")
	}
}

func TestWAFNotifier_Notify_ClientError(t *testing.T) {
	client := &mockWAFClient{err: errors.New("waf api error")}
	n := newWAFNotifierWithClient(client, "set-id", "set-name", types.ScopeRegional, "us-east-1")
	secret := newWAFSecret(3)
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client")
	}
}
