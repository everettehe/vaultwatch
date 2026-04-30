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
	sentInput *sesv2.SendEmailInput
	err       error
}

func (m *mockSESTemplateClient) SendEmail(_ context.Context, params *sesv2.SendEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	m.sentInput = params
	return &sesv2.SendEmailOutput{}, m.err
}

func newSESTemplateSecret(days int) *vault.Secret {
	return &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(time.Duration(days) * 24 * time.Hour),
	}
}

func TestNewSESTemplateNotifier_MissingFrom(t *testing.T) {
	_, err := newSESTemplateNotifierWithClient(&mockSESTemplateClient{}, "", "to@example.com", "MyTemplate")
	if err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestNewSESTemplateNotifier_MissingTo(t *testing.T) {
	_, err := newSESTemplateNotifierWithClient(&mockSESTemplateClient{}, "from@example.com", "", "MyTemplate")
	if err == nil {
		t.Fatal("expected error for missing to")
	}
}

func TestNewSESTemplateNotifier_MissingTemplate(t *testing.T) {
	_, err := newSESTemplateNotifierWithClient(&mockSESTemplateClient{}, "from@example.com", "to@example.com", "")
	if err == nil {
		t.Fatal("expected error for missing template name")
	}
}

func TestSESTemplateNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockSESTemplateClient{}
	n, err := newSESTemplateNotifierWithClient(mock, "from@example.com", "to@example.com", "VaultAlert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secret := newSESTemplateSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected notify error: %v", err)
	}
	if mock.sentInput == nil {
		t.Fatal("expected SendEmail to be called")
	}
	if *mock.sentInput.Content.Template.TemplateName != "VaultAlert" {
		t.Errorf("expected template name VaultAlert, got %s", *mock.sentInput.Content.Template.TemplateName)
	}
}

func TestSESTemplateNotifier_Notify_Expired(t *testing.T) {
	mock := &mockSESTemplateClient{}
	n, _ := newSESTemplateNotifierWithClient(mock, "from@example.com", "to@example.com", "VaultAlert")
	secret := newSESTemplateSecret(-1)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSESTemplateNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockSESTemplateClient{err: errors.New("send failed")}
	n, _ := newSESTemplateNotifierWithClient(mock, "from@example.com", "to@example.com", "VaultAlert")
	secret := newSESTemplateSecret(3)
	err := n.Notify(context.Background(), secret)
	if err == nil {
		t.Fatal("expected error from client")
	}
}
