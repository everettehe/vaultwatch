package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockEMRClient struct {
	called bool
	err    error
}

func (m *mockEMRClient) AddJobFlowSteps(_ context.Context, _ *emr.AddJobFlowStepsInput, _ ...func(*emr.Options)) (*emr.AddJobFlowStepsOutput, error) {
	m.called = true
	return &emr.AddJobFlowStepsOutput{}, m.err
}

func newEMRSecret() *vault.Secret {
	return &vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
}

func TestNewEMRNotifier_MissingClusterID(t *testing.T) {
	_, err := NewEMRNotifier("", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing cluster ID")
	}
}

func TestNewEMRNotifier_MissingRegion(t *testing.T) {
	_, err := NewEMRNotifier("j-XXXXXXXXXXXXX", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewEMRNotifier_Valid(t *testing.T) {
	client := &mockEMRClient{}
	n := newEMRNotifierWithClient(client, "j-XXXXXXXXXXXXX", "us-east-1")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestEMRNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockEMRClient{}
	n := newEMRNotifierWithClient(client, "j-XXXXXXXXXXXXX", "us-east-1")
	secret := newEMRSecret()

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Fatal("expected AddJobFlowSteps to be called")
	}
}

func TestEMRNotifier_Notify_Expired(t *testing.T) {
	client := &mockEMRClient{}
	n := newEMRNotifierWithClient(client, "j-XXXXXXXXXXXXX", "us-east-1")
	secret := &vault.Secret{
		Path:      "secret/api/key",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Fatal("expected AddJobFlowSteps to be called")
	}
}

func TestEMRNotifier_Notify_ClientError(t *testing.T) {
	client := &mockEMRClient{err: errors.New("aws error")}
	n := newEMRNotifierWithClient(client, "j-XXXXXXXXXXXXX", "us-east-1")
	secret := newEMRSecret()

	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client")
	}
}
