package notifier

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockAWSConfigClient struct {
	called bool
	err    error
}

func (m *mockAWSConfigClient) PutEvaluations(_ context.Context, _ *configservice.PutEvaluationsInput, _ ...func(*configservice.Options)) (*configservice.PutEvaluationsOutput, error) {
	m.called = true
	return &configservice.PutEvaluationsOutput{}, m.err
}

func newAWSConfigSecret(daysUntil int) *vault.Secret {
	return &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour),
	}
}

func TestNewAWSConfigNotifier_MissingResultToken(t *testing.T) {
	_, err := newAWSConfigNotifierWithClient(&mockAWSConfigClient{}, "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing result token")
	}
}

func TestNewAWSConfigNotifier_MissingRegion(t *testing.T) {
	_, err := newAWSConfigNotifierWithClient(&mockAWSConfigClient{}, "token-abc", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewAWSConfigNotifier_Valid(t *testing.T) {
	n, err := newAWSConfigNotifierWithClient(&mockAWSConfigClient{}, "token-abc", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestAWSConfigNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockAWSConfigClient{}
	n, _ := newAWSConfigNotifierWithClient(mock, "token-abc", "us-east-1")
	secret := newAWSConfigSecret(3)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected PutEvaluations to be called")
	}
}

func TestAWSConfigNotifier_Notify_Expired(t *testing.T) {
	mock := &mockAWSConfigClient{}
	n, _ := newAWSConfigNotifierWithClient(mock, "token-abc", "us-east-1")
	secret := newAWSConfigSecret(-1)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected PutEvaluations to be called")
	}
}

func TestAWSConfigNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockAWSConfigClient{err: fmt.Errorf("api error")}
	n, _ := newAWSConfigNotifierWithClient(mock, "token-abc", "us-east-1")
	secret := newAWSConfigSecret(5)
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client")
	}
}
