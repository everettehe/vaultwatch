package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/wakeward/vaultwatch/internal/vault"
)

type mockSQSClient struct {
	sentInput *sqs.SendMessageInput
	err       error
}

func (m *mockSQSClient) SendMessage(_ context.Context, params *sqs.SendMessageInput, _ ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	m.sentInput = params
	return &sqs.SendMessageOutput{}, m.err
}

func newSQSSecret(daysUntil int) *vault.Secret {
	return &vault.Secret{
		Path:      "secret/sqs/test",
		ExpiresAt: time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour),
	}
}

func TestNewSQSNotifier_MissingQueueURL(t *testing.T) {
	_, err := newSQSNotifierWithClient(&mockSQSClient{}, "")
	if err == nil {
		t.Fatal("expected error for missing queue URL")
	}
}

func TestNewSQSNotifier_Valid(t *testing.T) {
	n, err := newSQSNotifierWithClient(&mockSQSClient{}, "https://sqs.us-east-1.amazonaws.com/123456789/my-queue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSQSNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockSQSClient{}
	n, _ := newSQSNotifierWithClient(client, "https://sqs.us-east-1.amazonaws.com/123456789/my-queue")

	if err := n.Notify(context.Background(), newSQSSecret(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.sentInput == nil {
		t.Fatal("expected SendMessage to be called")
	}
	if client.sentInput.MessageBody == nil || *client.sentInput.MessageBody == "" {
		t.Fatal("expected non-empty message body")
	}
}

func TestSQSNotifier_Notify_Expired(t *testing.T) {
	client := &mockSQSClient{}
	n, _ := newSQSNotifierWithClient(client, "https://sqs.us-east-1.amazonaws.com/123456789/my-queue")

	if err := n.Notify(context.Background(), newSQSSecret(-1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.sentInput == nil {
		t.Fatal("expected SendMessage to be called")
	}
}

func TestSQSNotifier_Notify_ClientError(t *testing.T) {
	client := &mockSQSClient{err: errors.New("sqs unavailable")}
	n, _ := newSQSNotifierWithClient(client, "https://sqs.us-east-1.amazonaws.com/123456789/my-queue")

	err := n.Notify(context.Background(), newSQSSecret(3))
	if err == nil {
		t.Fatal("expected error from client")
	}
}

func TestSQSNotifier_ImplementsInterface(t *testing.T) {
	client := &mockSQSClient{}
	n, _ := newSQSNotifierWithClient(client, "https://sqs.us-east-1.amazonaws.com/123456789/my-queue")
	var _ Notifier = n
}
