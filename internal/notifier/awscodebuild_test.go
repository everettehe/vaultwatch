package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockCodeBuildClient struct {
	called bool
	err    error
}

func (m *mockCodeBuildClient) StartBuild(_ context.Context, _ *codebuild.StartBuildInput, _ ...func(*codebuild.Options)) (*codebuild.StartBuildOutput, error) {
	m.called = true
	return &codebuild.StartBuildOutput{}, m.err
}

func newCodeBuildSecret() *vault.Secret {
	return &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}
}

func TestNewCodeBuildNotifier_MissingProjectName(t *testing.T) {
	_, err := NewCodeBuildNotifier("", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing project name")
	}
}

func TestNewCodeBuildNotifier_MissingRegion(t *testing.T) {
	_, err := NewCodeBuildNotifier("my-project", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewCodeBuildNotifier_Valid(t *testing.T) {
	mock := &mockCodeBuildClient{}
	n := newCodeBuildNotifierWithClient(mock, "my-project", "us-east-1")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	if n.projectName != "my-project" {
		t.Errorf("expected project name 'my-project', got %q", n.projectName)
	}
}

func TestCodeBuildNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockCodeBuildClient{}
	n := newCodeBuildNotifierWithClient(mock, "rotation-project", "us-east-1")
	secret := newCodeBuildSecret()
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected StartBuild to be called")
	}
}

func TestCodeBuildNotifier_Notify_Expired(t *testing.T) {
	mock := &mockCodeBuildClient{}
	n := newCodeBuildNotifierWithClient(mock, "rotation-project", "us-east-1")
	secret := &vault.Secret{
		Path:      "secret/expired",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected StartBuild to be called")
	}
}

func TestCodeBuildNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockCodeBuildClient{err: errors.New("AWS error")}
	n := newCodeBuildNotifierWithClient(mock, "rotation-project", "us-east-1")
	secret := newCodeBuildSecret()
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client")
	}
}
