package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockIoTPublisher struct {
	calledWith *iotdataplane.PublishInput
	err        error
}

func (m *mockIoTPublisher) Publish(_ context.Context, params *iotdataplane.PublishInput, _ ...func(*iotdataplane.Options)) (*iotdataplane.PublishOutput, error) {
	m.calledWith = params
	return &iotdataplane.PublishOutput{}, m.err
}

func newIoTSecret(daysLeft int) vault.Secret {
	expiry := time.Now().Add(time.Duration(daysLeft) * 24 * time.Hour)
	return vault.Secret{Path: "secret/iot/device", ExpiresAt: expiry}
}

func TestNewAWSIoTNotifier_MissingTopic(t *testing.T) {
	_, err := newAWSIoTNotifierWithClient(&mockIoTPublisher{}, "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing topic")
	}
}

func TestNewAWSIoTNotifier_MissingRegion(t *testing.T) {
	_, err := newAWSIoTNotifierWithClient(&mockIoTPublisher{}, "vaultwatch/alerts", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewAWSIoTNotifier_Valid(t *testing.T) {
	n, err := newAWSIoTNotifierWithClient(&mockIoTPublisher{}, "vaultwatch/alerts", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestAWSIoTNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockIoTPublisher{}
	n, _ := newAWSIoTNotifierWithClient(mock, "vaultwatch/alerts", "us-east-1")
	if err := n.Notify(context.Background(), newIoTSecret(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calledWith == nil {
		t.Fatal("expected Publish to be called")
	}
	if *mock.calledWith.Topic != "vaultwatch/alerts" {
		t.Errorf("unexpected topic: %s", *mock.calledWith.Topic)
	}
	if len(mock.calledWith.Payload) == 0 {
		t.Error("expected non-empty payload")
	}
}

func TestAWSIoTNotifier_Notify_Expired(t *testing.T) {
	mock := &mockIoTPublisher{}
	n, _ := newAWSIoTNotifierWithClient(mock, "vaultwatch/alerts", "us-east-1")
	if err := n.Notify(context.Background(), newIoTSecret(-1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAWSIoTNotifier_Notify_PublishError(t *testing.T) {
	mock := &mockIoTPublisher{err: errors.New("iot publish failed")}
	n, _ := newAWSIoTNotifierWithClient(mock, "vaultwatch/alerts", "us-east-1")
	err := n.Notify(context.Background(), newIoTSecret(3))
	if err == nil {
		t.Fatal("expected error from publish failure")
	}
}
