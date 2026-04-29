package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ses"
)

type mockSESRawClient struct {
	called bool
	err    error
}

func (m *mockSESRawClient) SendRawEmail(_ context.Context, _ *ses.SendRawEmailInput, _ ...func(*ses.Options)) (*ses.SendRawEmailOutput, error) {
	m.called = true
	return &ses.SendRawEmailOutput{}, m.err
}

func newSESRawSecret() *Secret {
	expiry := time.Now().Add(48 * time.Hour)
	return &Secret{Path: "secret/ses-raw", ExpiresAt: &expiry}
}

func TestNewSESRawNotifier_MissingFrom(t *testing.T) {
	_, err := NewSESRawNotifier("", "to@example.com", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestNewSESRawNotifier_MissingTo(t *testing.T) {
	_, err := NewSESRawNotifier("from@example.com", "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing to")
	}
}

func TestSESRawNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockSESRawClient{}
	n := newSESRawNotifierWithClient(mock, "from@example.com", "to@example.com")
	if err := n.Notify(context.Background(), newSESRawSecret()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected SendRawEmail to be called")
	}
}

func TestSESRawNotifier_Notify_Expired(t *testing.T) {
	expiry := time.Now().Add(-1 * time.Hour)
	secret := &Secret{Path: "secret/ses-raw-expired", ExpiresAt: &expiry}
	mock := &mockSESRawClient{}
	n := newSESRawNotifierWithClient(mock, "from@example.com", "to@example.com")
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSESRawNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockSESRawClient{err: errors.New("send failed")}
	n := newSESRawNotifierWithClient(mock, "from@example.com", "to@example.com")
	if err := n.Notify(context.Background(), newSESRawSecret()); err == nil {
		t.Fatal("expected error from client failure")
	}
}
