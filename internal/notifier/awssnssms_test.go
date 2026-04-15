package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockSNSSMSClient struct {
	called bool
	input  *sns.PublishInput
	err    error
}

func (m *mockSNSSMSClient) Publish(_ context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.called = true
	m.input = params
	return &sns.PublishOutput{}, m.err
}

func newSNSSMSSecret(daysUntil int) *vault.Secret {
	return &vault.Secret{
		Path:      "secret/sms-test",
		ExpiresAt: time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour),
	}
}

func TestNewSNSSMSNotifier_MissingPhoneNumber(t *testing.T) {
	_, err := newSNSSMSNotifierWithClient(&mockSNSSMSClient{}, "", "")
	if err == nil {
		t.Fatal("expected error for missing phone number")
	}
}

func TestNewSNSSMSNotifier_Valid(t *testing.T) {
	n, err := newSNSSMSNotifierWithClient(&mockSNSSMSClient{}, "+15555550100", "VaultWatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.phoneNumber != "+15555550100" {
		t.Errorf("expected phone number +15555550100, got %s", n.phoneNumber)
	}
	if n.senderID != "VaultWatch" {
		t.Errorf("expected senderID VaultWatch, got %s", n.senderID)
	}
}

func TestSNSSMSNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockSNSSMSClient{}
	n, _ := newSNSSMSNotifierWithClient(mock, "+15555550100", "")
	secret := newSNSSMSSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected Publish to be called")
	}
	if mock.input == nil || mock.input.PhoneNumber == nil {
		t.Error("expected phone number to be set")
	}
}

func TestSNSSMSNotifier_Notify_Expired(t *testing.T) {
	mock := &mockSNSSMSClient{}
	n, _ := newSNSSMSNotifierWithClient(mock, "+15555550100", "")
	secret := newSNSSMSSecret(-1)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected Publish to be called")
	}
}

func TestSNSSMSNotifier_Notify_PublishError(t *testing.T) {
	mock := &mockSNSSMSClient{err: errors.New("sns error")}
	n, _ := newSNSSMSNotifierWithClient(mock, "+15555550100", "")
	secret := newSNSSMSSecret(3)
	err := n.Notify(context.Background(), secret)
	if err == nil {
		t.Fatal("expected error from failed publish")
	}
}

func TestSNSSMSNotifier_ImplementsInterface(t *testing.T) {
	n, _ := newSNSSMSNotifierWithClient(&mockSNSSMSClient{}, "+15555550100", "")
	var _ Notifier = n
}
