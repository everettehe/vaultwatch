package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/younsl/vaultwatch/internal/vault"
)

type mockCWEventsClient struct {
	err error
}

func (m *mockCWEventsClient) PutEvents(_ context.Context, _ *cloudwatchevents.PutEventsInput, _ ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error) {
	return &cloudwatchevents.PutEventsOutput{}, m.err
}

func newCWEventsSecret(days int) *vault.Secret {
	return &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(time.Duration(days) * 24 * time.Hour),
	}
}

func TestNewCloudWatchEventsNotifier_MissingEventBus(t *testing.T) {
	_, err := newCloudWatchEventsNotifierWithClient(&mockCWEventsClient{}, "", "vaultwatch", "VaultSecretExpiration")
	if err == nil {
		t.Fatal("expected error for missing event bus")
	}
}

func TestNewCloudWatchEventsNotifier_DefaultSourceAndDetailType(t *testing.T) {
	n, err := newCloudWatchEventsNotifierWithClient(&mockCWEventsClient{}, "my-event-bus", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.source != "vaultwatch" {
		t.Errorf("expected source 'vaultwatch', got %q", n.source)
	}
	if n.detailType != "VaultSecretExpiration" {
		t.Errorf("expected detailType 'VaultSecretExpiration', got %q", n.detailType)
	}
}

func TestNewCloudWatchEventsNotifier_Valid(t *testing.T) {
	n, err := newCloudWatchEventsNotifierWithClient(&mockCWEventsClient{}, "my-event-bus", "vaultwatch", "VaultSecretExpiration")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestCloudWatchEventsNotifier_Notify_ExpiringSoon(t *testing.T) {
	n, _ := newCloudWatchEventsNotifierWithClient(&mockCWEventsClient{}, "my-event-bus", "", "")
	secret := newCWEventsSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloudWatchEventsNotifier_Notify_Expired(t *testing.T) {
	n, _ := newCloudWatchEventsNotifierWithClient(&mockCWEventsClient{}, "my-event-bus", "", "")
	secret := newCWEventsSecret(-1)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloudWatchEventsNotifier_Notify_ClientError(t *testing.T) {
	n, _ := newCloudWatchEventsNotifierWithClient(&mockCWEventsClient{err: errors.New("put failed")}, "my-event-bus", "", "")
	secret := newCWEventsSecret(3)
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client")
	}
}

func TestCloudWatchEventsNotifier_ImplementsInterface(t *testing.T) {
	n, _ := newCloudWatchEventsNotifierWithClient(&mockCWEventsClient{}, "my-event-bus", "", "")
	var _ Notifier = n
}
