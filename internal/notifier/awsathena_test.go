package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/athena"

	vaultwatch "github.com/yourusername/vaultwatch/internal/vault"
)

type mockAthenaClient struct {
	called bool
	err    error
}

func (m *mockAthenaClient) StartQueryExecution(_ context.Context, _ *athena.StartQueryExecutionInput, _ ...func(*athena.Options)) (*athena.StartQueryExecutionOutput, error) {
	m.called = true
	return &athena.StartQueryExecutionOutput{}, m.err
}

func newAthenaSecret(daysUntil int) *vaultwatch.Secret {
	expiry := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return &vaultwatch.Secret{
		Path:      "secret/athena-test",
		ExpiresAt: expiry,
	}
}

func TestNewAthenaNotifier_MissingDatabase(t *testing.T) {
	_, err := NewAthenaNotifier("", "primary", "s3://bucket/", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing database")
	}
}

func TestNewAthenaNotifier_MissingOutputLocation(t *testing.T) {
	_, err := NewAthenaNotifier("mydb", "primary", "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing output_location")
	}
}

func TestNewAthenaNotifier_MissingRegion(t *testing.T) {
	_, err := NewAthenaNotifier("mydb", "primary", "s3://bucket/", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewAthenaNotifier_DefaultWorkgroup(t *testing.T) {
	client := &mockAthenaClient{}
	n := newAthenaNotifierWithClient(client, "mydb", "primary", "s3://bucket/")
	if n.workgroup != "primary" {
		t.Errorf("expected workgroup 'primary', got %q", n.workgroup)
	}
}

func TestAthenaNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockAthenaClient{}
	n := newAthenaNotifierWithClient(client, "mydb", "primary", "s3://bucket/")
	secret := newAthenaSecret(5)

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Error("expected StartQueryExecution to be called")
	}
}

func TestAthenaNotifier_Notify_Expired(t *testing.T) {
	client := &mockAthenaClient{}
	n := newAthenaNotifierWithClient(client, "mydb", "primary", "s3://bucket/")
	secret := newAthenaSecret(-1)

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Error("expected StartQueryExecution to be called")
	}
}

func TestAthenaNotifier_Notify_ClientError(t *testing.T) {
	client := &mockAthenaClient{err: errors.New("athena unavailable")}
	n := newAthenaNotifierWithClient(client, "mydb", "primary", "s3://bucket/")
	secret := newAthenaSecret(3)

	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client, got nil")
	}
}
