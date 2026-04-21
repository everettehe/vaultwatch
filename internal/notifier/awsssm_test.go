package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"

	vaultwatch "github.com/yourusername/vaultwatch/internal/vault"
)

type mockSSMClient struct {
	calledWith *ssm.PutParameterInput
	err        error
}

func (m *mockSSMClient) PutParameter(_ context.Context, params *ssm.PutParameterInput, _ ...func(*ssm.Options)) (*ssm.PutParameterOutput, error) {
	m.calledWith = params
	return &ssm.PutParameterOutput{}, m.err
}

func newSSMSecret() *vaultwatch.Secret {
	return &vaultwatch.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().UTC().Add(5 * 24 * time.Hour),
	}
}

func TestNewSSMNotifier_MissingParamPath(t *testing.T) {
	_, err := newSSMNotifierWithClient(&mockSSMClient{}, "")
	if err == nil {
		t.Fatal("expected error for empty param path")
	}
}

func TestNewSSMNotifier_Valid(t *testing.T) {
	n, err := newSSMNotifierWithClient(&mockSSMClient{}, "/vaultwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSSMNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockSSMClient{}
	n, _ := newSSMNotifierWithClient(mock, "/vaultwatch/alerts")
	secret := newSSMSecret()

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calledWith == nil {
		t.Fatal("expected PutParameter to be called")
	}
	if mock.calledWith.Value == nil || *mock.calledWith.Value == "" {
		t.Error("expected non-empty parameter value")
	}
}

func TestSSMNotifier_Notify_Expired(t *testing.T) {
	mock := &mockSSMClient{}
	n, _ := newSSMNotifierWithClient(mock, "/vaultwatch/alerts")
	secret := &vaultwatch.Secret{
		Path:      "secret/myapp/old-token",
		ExpiresAt: time.Now().UTC().Add(-24 * time.Hour),
	}

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSSMNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockSSMClient{err: errors.New("ssm unavailable")}
	n, _ := newSSMNotifierWithClient(mock, "/vaultwatch/alerts")

	err := n.Notify(context.Background(), newSSMSecret())
	if err == nil {
		t.Fatal("expected error from client")
	}
}

func TestSSMNotifier_ImplementsInterface(t *testing.T) {
	n, _ := newSSMNotifierWithClient(&mockSSMClient{}, "/vaultwatch")
	var _ Notifier = n
}
