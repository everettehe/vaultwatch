package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockRoute53Client struct {
	called bool
	err    error
}

func (m *mockRoute53Client) ChangeResourceRecordSets(_ context.Context, _ *route53.ChangeResourceRecordSetsInput, _ ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error) {
	m.called = true
	return &route53.ChangeResourceRecordSetsOutput{}, m.err
}

func newRoute53Secret(daysUntil int) *vault.Secret {
	expiry := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return &vault.Secret{Path: "secret/route53/test", ExpiresAt: expiry}
}

func TestNewRoute53Notifier_MissingZoneID(t *testing.T) {
	_, err := NewRoute53Notifier("", "vaultwatch.example.com.", 300)
	if err == nil {
		t.Fatal("expected error for missing zone ID")
	}
}

func TestNewRoute53Notifier_MissingRecord(t *testing.T) {
	_, err := NewRoute53Notifier("Z1234567890", "", 300)
	if err == nil {
		t.Fatal("expected error for missing record name")
	}
}

func TestNewRoute53Notifier_DefaultTTL(t *testing.T) {
	client := &mockRoute53Client{}
	n := newRoute53NotifierWithClient(client, "Z1234567890", "vaultwatch.example.com.", 0)
	if n.ttl != 0 {
		t.Errorf("expected ttl 0 from helper, got %d", n.ttl)
	}
}

func TestRoute53Notifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockRoute53Client{}
	n := newRoute53NotifierWithClient(client, "Z1234567890", "vaultwatch.example.com.", 300)
	s := newRoute53Secret(5)
	if err := n.Notify(context.Background(), s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Error("expected ChangeResourceRecordSets to be called")
	}
}

func TestRoute53Notifier_Notify_Expired(t *testing.T) {
	client := &mockRoute53Client{}
	n := newRoute53NotifierWithClient(client, "Z1234567890", "vaultwatch.example.com.", 300)
	s := newRoute53Secret(-1)
	if err := n.Notify(context.Background(), s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Error("expected ChangeResourceRecordSets to be called")
	}
}

func TestRoute53Notifier_Notify_ClientError(t *testing.T) {
	client := &mockRoute53Client{err: errors.New("route53 API error")}
	n := newRoute53NotifierWithClient(client, "Z1234567890", "vaultwatch.example.com.", 300)
	s := newRoute53Secret(3)
	if err := n.Notify(context.Background(), s); err == nil {
		t.Fatal("expected error from client")
	}
}

func TestRoute53Notifier_ImplementsInterface(t *testing.T) {
	client := &mockRoute53Client{}
	n := newRoute53NotifierWithClient(client, "Z1234567890", "vaultwatch.example.com.", 300)
	var _ Notifier = n
}
