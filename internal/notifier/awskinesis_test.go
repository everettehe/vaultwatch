package notifier_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"

	"github.com/yourusername/vaultwatch/internal/notifier"
	vaultsecret "github.com/yourusername/vaultwatch/internal/vault"
)

type mockKinesisClient struct {
	calledWith *kinesis.PutRecordInput
	err        error
}

func (m *mockKinesisClient) PutRecord(_ context.Context, params *kinesis.PutRecordInput, _ ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error) {
	m.calledWith = params
	return &kinesis.PutRecordOutput{}, m.err
}

func newKinesisSecret(daysUntil int) *vaultsecret.Secret {
	return &vaultsecret.Secret{
		Path:      "secret/kinesis-test",
		ExpiresAt: time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour),
	}
}

func TestNewKinesisNotifier_MissingStreamName(t *testing.T) {
	_, err := notifier.NewKinesisNotifier("", "")
	if err == nil {
		t.Fatal("expected error for missing stream name")
	}
}

func TestNewKinesisNotifier_DefaultPartitionKey(t *testing.T) {
	mock := &mockKinesisClient{}
	n := notifier.NewKinesisNotifierWithClient(mock, "my-stream", "")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestKinesisNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockKinesisClient{}
	n := notifier.NewKinesisNotifierWithClient(mock, "vault-events", "vaultwatch")
	secret := newKinesisSecret(5)

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calledWith == nil {
		t.Fatal("expected PutRecord to be called")
	}
	if *mock.calledWith.StreamName != "vault-events" {
		t.Errorf("expected stream name 'vault-events', got %q", *mock.calledWith.StreamName)
	}
	if *mock.calledWith.PartitionKey != "vaultwatch" {
		t.Errorf("expected partition key 'vaultwatch', got %q", *mock.calledWith.PartitionKey)
	}
}

func TestKinesisNotifier_Notify_Expired(t *testing.T) {
	mock := &mockKinesisClient{}
	n := notifier.NewKinesisNotifierWithClient(mock, "vault-events", "vaultwatch")
	secret := newKinesisSecret(-1)

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mock.calledWith.Data) == 0 {
		t.Error("expected non-empty data in PutRecord call")
	}
}

func TestKinesisNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockKinesisClient{err: errors.New("kinesis unavailable")}
	n := notifier.NewKinesisNotifierWithClient(mock, "vault-events", "vaultwatch")
	secret := newKinesisSecret(3)

	err := n.Notify(context.Background(), secret)
	if err == nil {
		t.Fatal("expected error from client")
	}
}

func TestKinesisNotifier_ImplementsInterface(t *testing.T) {
	mock := &mockKinesisClient{}
	var _ notifier.Notifier = notifier.NewKinesisNotifierWithClient(mock, "stream", "key")
}
