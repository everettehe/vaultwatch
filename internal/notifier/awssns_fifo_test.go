package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockSNSFIFOClient struct {
	publishFn func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

func (m *mockSNSFIFOClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	return m.publishFn(ctx, params, optFns...)
}

func newFIFOSecret(days int) vault.Secret {
	return vault.Secret{
		Path:      "secret/fifo/test",
		ExpiresAt: time.Now().Add(time.Duration(days) * 24 * time.Hour),
	}
}

func TestNewSNSFIFONotifier_MissingTopicARN(t *testing.T) {
	_, err := newSNSFIFONotifierWithClient(&mockSNSFIFOClient{}, "", "group1")
	if err == nil {
		t.Fatal("expected error for missing topic ARN")
	}
}

func TestNewSNSFIFONotifier_MissingGroupID(t *testing.T) {
	_, err := newSNSFIFONotifierWithClient(&mockSNSFIFOClient{}, "arn:aws:sns:us-east-1:123456789012:alerts.fifo", "")
	if err == nil {
		t.Fatal("expected error for missing group ID")
	}
}

func TestNewSNSFIFONotifier_Valid(t *testing.T) {
	n, err := newSNSFIFONotifierWithClient(&mockSNSFIFOClient{}, "arn:aws:sns:us-east-1:123456789012:alerts.fifo", "vaultwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSNSFIFONotifier_Notify_ExpiringSoon(t *testing.T) {
	var capturedInput *sns.PublishInput
	client := &mockSNSFIFOClient{
		publishFn: func(ctx context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
			capturedInput = params
			return &sns.PublishOutput{}, nil
		},
	}
	n, _ := newSNSFIFONotifierWithClient(client, "arn:aws:sns:us-east-1:123456789012:alerts.fifo", "vaultwatch")
	if err := n.Notify(context.Background(), newFIFOSecret(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedInput == nil {
		t.Fatal("expected publish to be called")
	}
	if *capturedInput.MessageGroupId != "vaultwatch" {
		t.Errorf("expected group ID 'vaultwatch', got %q", *capturedInput.MessageGroupId)
	}
}

func TestSNSFIFONotifier_Notify_PublishError(t *testing.T) {
	client := &mockSNSFIFOClient{
		publishFn: func(_ context.Context, _ *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
			return nil, errors.New("publish failed")
		},
	}
	n, _ := newSNSFIFONotifierWithClient(client, "arn:aws:sns:us-east-1:123456789012:alerts.fifo", "vaultwatch")
	if err := n.Notify(context.Background(), newFIFOSecret(3)); err == nil {
		t.Fatal("expected error from publish failure")
	}
}
