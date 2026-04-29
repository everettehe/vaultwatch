package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockSESTemplateClient struct {
	called bool
	err    error
}

func (m *mockSESTemplateClient) SendEmail(_ context.Context, _ *sesv2.SendEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	m.called = true
	return &sesv2.SendEmailOutput{}, m.err
}

func newSESTemplateSecret() vault.Secret {
	return vault.Secret{
		Path:      "secret/myapp/api-key",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
}

func TestNewSESTemplateNotifier_MissingFrom(t *testing.T) {
	_, err := NewSESTemplateNotifier("", "to@example.com", "MyTemplate", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestNewSESTemplateNotifier_MissingTo(t *testing.T) {
	_, err := NewSESTemplateNotifier("from@example.com", "", "MyTemplate", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing to")
	}
}

func TestNewSESTemplateNotifier_MissingTemplate(t *testing.T) {
	_, err := NewSESTemplateNotifier("from@example.com", "to@example.com", "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestNewSESTemplateNotifier_MissingRegion(t *testing.T) {
	_, err := NewSESTemplateNotifier("from@example.com", "to@example.com", "MyTemplate", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestSESTemplateNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockSESTemplateClient{}
	n := newSESTemplateNotifierWithClient(client, "from@example.com", "to@example.com", "MyTemplate", "us-east-1")
	err := n.Notify(context.Background(), newSESTemplateSecret())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Fatal("expected SendEmail to be called")
	}
}

func TestSESTemplateNotifier_Notify_Error(t *testing.T) {
	client := &mockSESTemplateClient{err: errors.New("aws error")}
	n := newSESTemplateNotifierWithClient(client, "from@example.com", "to@example.com", "MyTemplate", "us-east-1")
	err := n.Notify(context.Background(), newSESTemplateSecret())
	if err == nil {
		t.Fatal("expected error from SendEmail")
	}
}
