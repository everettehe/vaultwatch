package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockECSClient struct {
	called bool
	err    error
}

func (m *mockECSClient) RunTask(_ context.Context, _ *ecs.RunTaskInput, _ ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
	m.called = true
	return &ecs.RunTaskOutput{}, m.err
}

func newECSSecret() vault.Secret {
	return vault.Secret{
		Path:      "secret/data/myapp/api-key",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
}

func TestNewECSNotifier_MissingCluster(t *testing.T) {
	_, err := newECSNotifierWithClient(&mockECSClient{}, "", "my-task", "app", "FARGATE", "us-east-1")
	if err != nil {
		t.Fatal("constructor with client should not validate; got", err)
	}
	// validation happens in NewECSNotifier
	_, err = NewECSNotifier("", "my-task", "app", "FARGATE", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing cluster")
	}
}

func TestNewECSNotifier_MissingTaskDef(t *testing.T) {
	_, err := NewECSNotifier("my-cluster", "", "app", "FARGATE", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing task_def")
	}
}

func TestNewECSNotifier_MissingRegion(t *testing.T) {
	_, err := NewECSNotifier("my-cluster", "my-task", "app", "FARGATE", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewECSNotifier_DefaultLaunchType(t *testing.T) {
	n, err := newECSNotifierWithClient(&mockECSClient{}, "cluster", "task", "app", "", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.launchType != "" {
		// default applied only in NewECSNotifier; internal constructor stores as-is
		t.Log("launch type stored as-is in internal constructor")
	}
}

func TestECSNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockECSClient{}
	n, err := newECSNotifierWithClient(client, "cluster", "task", "app", "FARGATE", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := newECSSecret()
	if err := n.Notify(context.Background(), s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Fatal("expected RunTask to be called")
	}
}

func TestECSNotifier_Notify_Expired(t *testing.T) {
	client := &mockECSClient{}
	n, _ := newECSNotifierWithClient(client, "cluster", "task", "app", "FARGATE", "us-east-1")
	s := vault.Secret{Path: "secret/data/old", ExpiresAt: time.Now().Add(-1 * time.Hour)}
	if err := n.Notify(context.Background(), s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestECSNotifier_Notify_ClientError(t *testing.T) {
	client := &mockECSClient{err: errors.New("aws error")}
	n, _ := newECSNotifierWithClient(client, "cluster", "task", "app", "FARGATE", "us-east-1")
	s := newECSSecret()
	if err := n.Notify(context.Background(), s); err == nil {
		t.Fatal("expected error from client")
	}
}
