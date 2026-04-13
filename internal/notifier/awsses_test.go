package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockSESClient struct {
	called bool
	input  *ses.SendEmailInput
	err    error
}

func (m *mockSESClient) SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
	m.called = true
	m.input = params
	return &ses.SendEmailOutput{}, m.err
}

func newSESSecret(daysUntil int) vault.Secret {
	expiry := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return vault.Secret{Path: "secret/ses-test", ExpiresAt: &expiry}
}

func TestNewSESNotifier_MissingFrom(t *testing.T) {
	_, err := NewSESNotifier("", []string{"to@example.com"}, "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing from address")
	}
}

func TestNewSESNotifier_MissingTo(t *testing.T) {
	_, err := NewSESNotifier("from@example.com", []string{}, "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing recipients")
	}
}

func TestSESNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockSESClient{}
	n := newSESNotifierWithClient(mock, "from@example.com", []string{"to@example.com"})

	secret := newSESSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Fatal("expected SendEmail to be called")
	}
	if mock.input == nil || mock.input.Source == nil || *mock.input.Source != "from@example.com" {
		t.Errorf("unexpected source address")
	}
}

func TestSESNotifier_Notify_Expired(t *testing.T) {
	mock := &mockSESClient{}
	n := newSESNotifierWithClient(mock, "from@example.com", []string{"ops@example.com"})

	secret := newSESSecret(-1)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Fatal("expected SendEmail to be called")
	}
}

func TestSESNotifier_Notify_SendError(t *testing.T) {
	mock := &mockSESClient{err: errors.New("SES throttled")}
	n := newSESNotifierWithClient(mock, "from@example.com", []string{"to@example.com"})

	secret := newSESSecret(3)
	err := n.Notify(context.Background(), secret)
	if err == nil {
		t.Fatal("expected error from SendEmail failure")
	}
}

func TestSESNotifier_ImplementsInterface(t *testing.T) {
	mock := &mockSESClient{}
	var _ Notifier = newSESNotifierWithClient(mock, "from@example.com", []string{"to@example.com"})
}
