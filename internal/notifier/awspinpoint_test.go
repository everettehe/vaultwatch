package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockPinpointSender struct {
	err error
}

func (m *mockPinpointSender) SendMessages(_ context.Context, _ *pinpoint.SendMessagesInput, _ ...func(*pinpoint.Options)) (*pinpoint.SendMessagesOutput, error) {
	return &pinpoint.SendMessagesOutput{}, m.err
}

func newPinpointSecret(daysUntil int) *vault.Secret {
	return &vault.Secret{
		Path:      "secret/pinpoint/test",
		ExpiresAt: time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour),
	}
}

func TestNewPinpointNotifier_MissingAppID(t *testing.T) {
	_, err := newPinpointNotifierWithClient(&mockPinpointSender{}, "", "+10000000000", "+19999999999", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing app_id")
	}
}

func TestNewPinpointNotifier_MissingDestNumber(t *testing.T) {
	_, err := newPinpointNotifierWithClient(&mockPinpointSender{}, "app-123", "+10000000000", "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing dest_number")
	}
}

func TestNewPinpointNotifier_MissingRegion(t *testing.T) {
	_, err := newPinpointNotifierWithClient(&mockPinpointSender{}, "app-123", "+10000000000", "+19999999999", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewPinpointNotifier_Valid(t *testing.T) {
	n, err := newPinpointNotifierWithClient(&mockPinpointSender{}, "app-123", "+10000000000", "+19999999999", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestPinpointNotifier_Notify_ExpiringSoon(t *testing.T) {
	n, _ := newPinpointNotifierWithClient(&mockPinpointSender{}, "app-123", "+10000000000", "+19999999999", "us-east-1")
	if err := n.Notify(context.Background(), newPinpointSecret(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPinpointNotifier_Notify_Expired(t *testing.T) {
	n, _ := newPinpointNotifierWithClient(&mockPinpointSender{}, "app-123", "+10000000000", "+19999999999", "us-east-1")
	if err := n.Notify(context.Background(), newPinpointSecret(-1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPinpointNotifier_Notify_SendError(t *testing.T) {
	n, _ := newPinpointNotifierWithClient(&mockPinpointSender{err: errors.New("send failed")}, "app-123", "+10000000000", "+19999999999", "us-east-1")
	if err := n.Notify(context.Background(), newPinpointSecret(3)); err == nil {
		t.Fatal("expected error from failed send")
	}
}
