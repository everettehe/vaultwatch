package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockSESBulkClient struct {
	err error
}

func (m *mockSESBulkClient) SendBulkEmail(_ context.Context, _ *sesv2.SendBulkEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendBulkEmailOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &sesv2.SendBulkEmailOutput{}, nil
}

func newSESBulkSecret() *vault.Secret {
	return &vault.Secret{
		Path:      "secret/myapp/api-key",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}
}

func TestNewSESSendBulkNotifier_MissingFrom(t *testing.T) {
	_, err := NewSESSendBulkNotifier("", []string{"to@example.com"}, "MyTemplate", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing from address")
	}
}

func TestNewSESSendBulkNotifier_MissingTo(t *testing.T) {
	_, err := NewSESSendBulkNotifier("from@example.com", nil, "MyTemplate", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing to addresses")
	}
}

func TestNewSESSendBulkNotifier_MissingTemplate(t *testing.T) {
	_, err := NewSESSendBulkNotifier("from@example.com", []string{"to@example.com"}, "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing template name")
	}
}

func TestNewSESSendBulkNotifier_MissingRegion(t *testing.T) {
	_, err := NewSESSendBulkNotifier("from@example.com", []string{"to@example.com"}, "MyTemplate", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestSESSendBulkNotifier_Notify_ExpiringSoon(t *testing.T) {
	n := newSESSendBulkNotifierWithClient(
		&mockSESBulkClient{},
		"from@example.com",
		[]string{"a@example.com", "b@example.com"},
		"ExpiryTemplate",
	)
	if err := n.Notify(context.Background(), newSESBulkSecret()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSESSendBulkNotifier_Notify_Error(t *testing.T) {
	n := newSESSendBulkNotifierWithClient(
		&mockSESBulkClient{err: errors.New("send failure")},
		"from@example.com",
		[]string{"to@example.com"},
		"ExpiryTemplate",
	)
	if err := n.Notify(context.Background(), newSESBulkSecret()); err == nil {
		t.Fatal("expected error from send failure")
	}
}

func TestSESSendBulkNotifier_ImplementsInterface(t *testing.T) {
	n := newSESSendBulkNotifierWithClient(
		&mockSESBulkClient{},
		"from@example.com",
		[]string{"to@example.com"},
		"ExpiryTemplate",
	)
	var _ Notifier = n
}
