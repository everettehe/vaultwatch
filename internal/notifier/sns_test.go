package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// mockSNSPublisher is a test double for snsPublisher.
type mockSNSPublisher struct {
	publishErr error
	captured   *sns.PublishInput
}

func (m *mockSNSPublisher) Publish(_ context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.captured = params
	if m.publishErr != nil {
		return nil, m.publishErr
	}
	return &sns.PublishOutput{}, nil
}

func TestNewSNSNotifier_MissingTopicARN(t *testing.T) {
	_, err := newSNSNotifierWithClient("", &mockSNSPublisher{})
	if err == nil {
		t.Fatal("expected error for empty topic ARN, got nil")
	}
}

func TestNewSNSNotifier_Valid(t *testing.T) {
	n, err := newSNSNotifierWithClient("arn:aws:sns:us-east-1:123456789012:alerts", &mockSNSPublisher{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSNSNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockSNSPublisher{}
	n, _ := newSNSNotifierWithClient("arn:aws:sns:us-east-1:123456789012:alerts", mock)

	secret := &vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.captured == nil {
		t.Fatal("expected Publish to be called")
	}
	if *mock.captured.TopicArn != "arn:aws:sns:us-east-1:123456789012:alerts" {
		t.Errorf("unexpected topic ARN: %s", *mock.captured.TopicArn)
	}
	if *mock.captured.Subject == "" {
		t.Error("expected non-empty subject")
	}
}

func TestSNSNotifier_Notify_Expired(t *testing.T) {
	mock := &mockSNSPublisher{}
	n, _ := newSNSNotifierWithClient("arn:aws:sns:us-east-1:123456789012:alerts", mock)

	secret := &vault.Secret{
		Path:      "secret/api/key",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSNSNotifier_Notify_PublishError(t *testing.T) {
	mock := &mockSNSPublisher{publishErr: errors.New("aws error")}
	n, _ := newSNSNotifierWithClient("arn:aws:sns:us-east-1:123456789012:alerts", mock)

	secret := &vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}

	err := n.Notify(secret)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
