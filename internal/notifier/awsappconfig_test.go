package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/appconfigdata"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockAppConfigClient struct {
	err error
}

func (m *mockAppConfigClient) StartConfigurationSession(_ context.Context, _ *appconfigdata.StartConfigurationSessionInput, _ ...func(*appconfigdata.Options)) (*appconfigdata.StartConfigurationSessionOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &appconfigdata.StartConfigurationSessionOutput{}, nil
}

func newAppConfigSecret(daysUntil int) *vault.Secret {
	return &vault.Secret{
		Path:      "secret/appconfig/test",
		ExpiresAt: time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour),
	}
}

func TestNewAppConfigNotifier_MissingApplication(t *testing.T) {
	_, err := newAppConfigNotifierWithClient(&mockAppConfigClient{}, "", "prod", "rotation-flag", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing application")
	}
}

func TestNewAppConfigNotifier_MissingEnvironment(t *testing.T) {
	_, err := newAppConfigNotifierWithClient(&mockAppConfigClient{}, "myapp", "", "rotation-flag", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestNewAppConfigNotifier_MissingProfile(t *testing.T) {
	_, err := newAppConfigNotifierWithClient(&mockAppConfigClient{}, "myapp", "prod", "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestNewAppConfigNotifier_Valid(t *testing.T) {
	n, err := newAppConfigNotifierWithClient(&mockAppConfigClient{}, "myapp", "prod", "rotation-flag", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestAppConfigNotifier_Notify_ExpiringSoon(t *testing.T) {
	n, _ := newAppConfigNotifierWithClient(&mockAppConfigClient{}, "myapp", "prod", "rotation-flag", "us-east-1")
	secret := newAppConfigSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAppConfigNotifier_Notify_Expired(t *testing.T) {
	n, _ := newAppConfigNotifierWithClient(&mockAppConfigClient{}, "myapp", "prod", "rotation-flag", "us-east-1")
	secret := newAppConfigSecret(-1)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAppConfigNotifier_Notify_ClientError(t *testing.T) {
	n, _ := newAppConfigNotifierWithClient(&mockAppConfigClient{err: errors.New("aws error")}, "myapp", "prod", "rotation-flag", "us-east-1")
	secret := newAppConfigSecret(3)
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client")
	}
}
