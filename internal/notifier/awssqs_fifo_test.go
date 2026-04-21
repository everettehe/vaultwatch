package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	vaultsecret "github.com/yourusername/vaultwatch/internal/vault"
)

type mockSQSFIFOClient struct {
	sentInput *sqs.SendMessageInput
	err       error
}

func (m *mockSQSFIFOClient) SendMessage(_ context.Context, params *sqs.SendMessageInput, _ ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	m.sentInput = params
	return &sqs.SendMessageOutput{}, m.err
}

func newSQSFIFOSecret() *vaultsecret.Secret {
	return &vaultsecret.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
}

func TestNewSQSFIFONotifier_MissingQueueURL(t *testing.T) {
	_, err := newSQSFIFONotifierWithClient(&mockSQSFIFOClient{}, "", "group1")
	if err == nil {
		t.Fatal("expected error for missing queue URL")
	}
}

func TestNewSQSFIFONotifier_MissingMessageGroup(t *testing.T) {
	_, err := newSQSFIFONotifierWithClient(&mockSQSFIFOClient{}, "https://sqs.us-east-1.amazonaws.com/123/queue.fifo", "")
	if err == nil {
		t.Fatal("expected error for missing message group ID")
	}
}

func TestNewSQSFIFONotifier_Valid(t *testing.T) {
	n, err := newSQSFIFONotifierWithClient(&mockSQSFIFOClient{}, "https://sqs.us-east-1.amazonaws.com/123/queue.fifo", "vaultwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSQSFIFONotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockSQSFIFOClient{}
	n, _ := newSQSFIFONotifierWithClient(client, "https://sqs.us-east-1.amazonaws.com/123/queue.fifo", "vaultwatch")
	secret := newSQSFIFOSecret()

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.sentInput == nil {
		t.Fatal("expected SendMessage to be called")
	}
	if *client.sentInput.MessageGroupId != "vaultwatch" {
		t.Errorf("expected message group 'vaultwatch', got %q", *client.sentInput.MessageGroupId)
	}
	if *client.sentInput.MessageDeduplicationId == "" {
		t.Error("expected non-empty deduplication ID")
	}
}

func TestSQSFIFONotifier_Notify_Expired(t *testing.T) {
	client := &mockSQSFIFOClient{}
	n, _ := newSQSFIFONotifierWithClient(client, "https://sqs.us-east-1.amazonaws.com/123/queue.fifo", "vaultwatch")
	secret := &vaultsecret.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSQSFIFONotifier_Notify_ClientError(t *testing.T) {
	client := &mockSQSFIFOClient{err: errors.New("sqs error")}
	n, _ := newSQSFIFONotifierWithClient(client, "https://sqs.us-east-1.amazonaws.com/123/queue.fifo", "vaultwatch")

	err := n.Notify(context.Background(), newSQSFIFOSecret())
	if err == nil {
		t.Fatal("expected error from client")
	}
}
