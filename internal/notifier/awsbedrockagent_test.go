package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockBedrockAgentClient struct {
	invokeErr error
	called    bool
}

func (m *mockBedrockAgentClient) InvokeAgent(_ context.Context, _ *bedrockagentruntime.InvokeAgentInput, _ ...func(*bedrockagentruntime.Options)) (*bedrockagentruntime.InvokeAgentOutput, error) {
	m.called = true
	if m.invokeErr != nil {
		return nil, m.invokeErr
	}
	return &bedrockagentruntime.InvokeAgentOutput{}, nil
}

func newBedrockSecret(daysLeft int) vault.Secret {
	return vault.Secret{
		Path:      "secret/bedrock/test",
		ExpiresAt: time.Now().Add(time.Duration(daysLeft) * 24 * time.Hour),
	}
}

func TestNewBedrockAgentNotifier_MissingAgentID(t *testing.T) {
	_, err := NewBedrockAgentNotifier("", "alias-123", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing agentID")
	}
}

func TestNewBedrockAgentNotifier_MissingAliasID(t *testing.T) {
	_, err := NewBedrockAgentNotifier("agent-123", "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing aliasID")
	}
}

func TestNewBedrockAgentNotifier_MissingRegion(t *testing.T) {
	_, err := NewBedrockAgentNotifier("agent-123", "alias-123", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewBedrockAgentNotifier_Valid(t *testing.T) {
	mock := &mockBedrockAgentClient{}
	n := newBedrockAgentNotifierWithClient(mock, "agent-123", "alias-123", "us-east-1")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	if n.agentID != "agent-123" {
		t.Errorf("expected agentID agent-123, got %s", n.agentID)
	}
}

func TestBedrockAgentNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockBedrockAgentClient{}
	n := newBedrockAgentNotifierWithClient(mock, "agent-123", "alias-123", "us-east-1")
	secret := newBedrockSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected InvokeAgent to be called")
	}
}

func TestBedrockAgentNotifier_Notify_Expired(t *testing.T) {
	mock := &mockBedrockAgentClient{}
	n := newBedrockAgentNotifierWithClient(mock, "agent-123", "alias-123", "us-east-1")
	secret := newBedrockSecret(-1)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBedrockAgentNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockBedrockAgentClient{invokeErr: errors.New("invoke failed")}
	n := newBedrockAgentNotifierWithClient(mock, "agent-123", "alias-123", "us-east-1")
	secret := newBedrockSecret(3)
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client")
	}
}
