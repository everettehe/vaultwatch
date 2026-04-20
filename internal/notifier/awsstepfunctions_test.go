package notifier

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type mockSFNClient struct {
	called bool
	lastInput *sfn.StartExecutionInput
	err    error
}

func (m *mockSFNClient) StartExecution(_ context.Context, params *sfn.StartExecutionInput, _ ...func(*sfn.Options)) (*sfn.StartExecutionOutput, error) {
	m.called = true
	m.lastInput = params
	return &sfn.StartExecutionOutput{}, m.err
}

func newSFNSecret(days int, expired bool) ExpiringSecret {
	return ExpiringSecret{
		Path:                "secret/db/password",
		DaysUntilExpiration: days,
		ExpiresAt:           time.Now().Add(time.Duration(days) * 24 * time.Hour),
		IsExpired:           expired,
	}
}

func TestNewStepFunctionsNotifier_MissingARN(t *testing.T) {
	_, err := newStepFunctionsNotifierWithClient(&mockSFNClient{}, "")
	if err == nil {
		t.Fatal("expected error for missing ARN")
	}
}

func TestNewStepFunctionsNotifier_Valid(t *testing.T) {
	n, err := newStepFunctionsNotifierWithClient(&mockSFNClient{}, "arn:aws:states:us-east-1:123456789012:stateMachine:MyMachine")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestStepFunctionsNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockSFNClient{}
	n, _ := newStepFunctionsNotifierWithClient(client, "arn:aws:states:us-east-1:123456789012:stateMachine:MyMachine")
	secret := newSFNSecret(5, false)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Fatal("expected StartExecution to be called")
	}
	if client.lastInput == nil || client.lastInput.Input == nil {
		t.Fatal("expected non-nil execution input")
	}
}

func TestStepFunctionsNotifier_Notify_Expired(t *testing.T) {
	client := &mockSFNClient{}
	n, _ := newStepFunctionsNotifierWithClient(client, "arn:aws:states:us-east-1:123456789012:stateMachine:MyMachine")
	secret := newSFNSecret(0, true)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Fatal("expected StartExecution to be called")
	}
}

func TestStepFunctionsNotifier_Notify_ClientError(t *testing.T) {
	client := &mockSFNClient{err: fmt.Errorf("throttled")}
	n, _ := newStepFunctionsNotifierWithClient(client, "arn:aws:states:us-east-1:123456789012:stateMachine:MyMachine")
	secret := newSFNSecret(3, false)
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from client")
	}
}

func TestStepFunctionsNotifier_ImplementsInterface(t *testing.T) {
	client := &mockSFNClient{}
	n, _ := newStepFunctionsNotifierWithClient(client, "arn:aws:states:us-east-1:123456789012:stateMachine:MyMachine")
	var _ Notifier = n
}
