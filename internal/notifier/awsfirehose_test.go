package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockFirehoseClient struct {
	called bool
	data   []byte
	err    error
}

func (m *mockFirehoseClient) PutRecord(_ context.Context, params *firehose.PutRecordInput, _ ...func(*firehose.Options)) (*firehose.PutRecordOutput, error) {
	m.called = true
	m.data = params.Record.Data
	return &firehose.PutRecordOutput{}, m.err
}

func newFirehoseSecret(daysLeft int) vault.Secret {
	return vault.Secret{
		Path:      "secret/firehose/test",
		ExpiresAt: time.Now().Add(time.Duration(daysLeft) * 24 * time.Hour),
	}
}

func TestNewFirehoseNotifier_MissingDeliveryStream(t *testing.T) {
	_, err := newFirehoseNotifierWithClient(&mockFirehoseClient{}, "")
	if err == nil {
		t.Fatal("expected error for empty delivery stream")
	}
}

func TestNewFirehoseNotifier_Valid(t *testing.T) {
	n, err := newFirehoseNotifierWithClient(&mockFirehoseClient{}, "my-stream")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestFirehoseNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockFirehoseClient{}
	n, _ := newFirehoseNotifierWithClient(client, "vault-alerts")

	secret := newFirehoseSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Error("expected PutRecord to be called")
	}
	if len(client.data) == 0 {
		t.Error("expected non-empty record data")
	}
	// Records must end with newline for line-delimited JSON.
	if client.data[len(client.data)-1] != '\n' {
		t.Error("expected record data to end with newline")
	}
}

func TestFirehoseNotifier_Notify_Expired(t *testing.T) {
	client := &mockFirehoseClient{}
	n, _ := newFirehoseNotifierWithClient(client, "vault-alerts")

	secret := newFirehoseSecret(-3)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Error("expected PutRecord to be called")
	}
}

func TestFirehoseNotifier_Notify_ClientError(t *testing.T) {
	client := &mockFirehoseClient{err: errors.New("firehose unavailable")}
	n, _ := newFirehoseNotifierWithClient(client, "vault-alerts")

	secret := newFirehoseSecret(10)
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client failure")
	}
}
