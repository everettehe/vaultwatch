package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockLambdaClient struct {
	invokedWith *lambda.InvokeInput
	err         error
}

func (m *mockLambdaClient) Invoke(_ context.Context, params *lambda.InvokeInput, _ ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
	m.invokedWith = params
	return &lambda.InvokeOutput{}, m.err
}

func newLambdaSecret() *vault.Secret {
	return &vault.Secret{
		Path:      "secret/myapp/api-key",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}
}

func TestNewLambdaNotifier_MissingFunctionName(t *testing.T) {
	_, err := newLambdaNotifierWithClient(&mockLambdaClient{}, "")
	if err == nil {
		t.Fatal("expected error for missing function name")
	}
}

func TestNewLambdaNotifier_Valid(t *testing.T) {
	n, err := newLambdaNotifierWithClient(&mockLambdaClient{}, "my-function")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestLambdaNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockLambdaClient{}
	n, _ := newLambdaNotifierWithClient(mock, "my-function")
	secret := newLambdaSecret()

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.invokedWith == nil {
		t.Fatal("expected Lambda to be invoked")
	}
	if *mock.invokedWith.FunctionName != "my-function" {
		t.Errorf("expected function name 'my-function', got %q", *mock.invokedWith.FunctionName)
	}
	if len(mock.invokedWith.Payload) == 0 {
		t.Error("expected non-empty payload")
	}
}

func TestLambdaNotifier_Notify_Expired(t *testing.T) {
	mock := &mockLambdaClient{}
	n, _ := newLambdaNotifierWithClient(mock, "my-function")
	secret := &vault.Secret{
		Path:      "secret/myapp/db-pass",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLambdaNotifier_Notify_InvocationError(t *testing.T) {
	mock := &mockLambdaClient{err: errors.New("connection refused")}
	n, _ := newLambdaNotifierWithClient(mock, "my-function")

	err := n.Notify(context.Background(), newLambdaSecret())
	if err == nil {
		t.Fatal("expected error from failed invocation")
	}
}

func TestLambdaNotifier_ImplementsInterface(t *testing.T) {
	n, _ := newLambdaNotifierWithClient(&mockLambdaClient{}, "fn")
	var _ Notifier = n
}
