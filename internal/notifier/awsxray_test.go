package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/xray"
	"github.com/your-org/vaultwatch/internal/vault"
)

type mockXRayClient struct {
	called bool
	err    error
	lastInput *xray.PutTraceSegmentsInput
}

func (m *mockXRayClient) PutTraceSegments(ctx context.Context, params *xray.PutTraceSegmentsInput, optFns ...func(*xray.Options)) (*xray.PutTraceSegmentsOutput, error) {
	m.called = true
	m.lastInput = params
	return &xray.PutTraceSegmentsOutput{}, m.err
}

func newXRaySecret(daysUntil int) *vault.Secret {
	expiry := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return &vault.Secret{Path: "secret/xray-test", Expiration: expiry}
}

func TestNewXRayNotifier_MissingRegion(t *testing.T) {
	_, err := NewXRayNotifier("", "vaultwatch")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewXRayNotifier_DefaultService(t *testing.T) {
	mock := &mockXRayClient{}
	n := newXRayNotifierWithClient(mock, "us-east-1", "vaultwatch")
	if n.service != "vaultwatch" {
		t.Errorf("expected service 'vaultwatch', got %q", n.service)
	}
}

func TestXRayNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockXRayClient{}
	n := newXRayNotifierWithClient(mock, "us-east-1", "vaultwatch")
	secret := newXRaySecret(5)

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Fatal("expected PutTraceSegments to be called")
	}
	if len(mock.lastInput.TraceSegmentDocuments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(mock.lastInput.TraceSegmentDocuments))
	}
}

func TestXRayNotifier_Notify_Expired(t *testing.T) {
	mock := &mockXRayClient{}
	n := newXRayNotifierWithClient(mock, "us-east-1", "vaultwatch")
	secret := newXRaySecret(-1)

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Fatal("expected PutTraceSegments to be called")
	}
}

func TestXRayNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockXRayClient{err: errors.New("xray unavailable")}
	n := newXRayNotifierWithClient(mock, "us-east-1", "vaultwatch")
	secret := newXRaySecret(3)

	err := n.Notify(context.Background(), secret)
	if err == nil {
		t.Fatal("expected error from client")
	}
}
