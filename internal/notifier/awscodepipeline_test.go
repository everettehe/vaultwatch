package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	vaultconfig "github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockCodePipelineClient struct {
	err error
}

func (m *mockCodePipelineClient) PutJobFailureResult(_ context.Context, _ *codepipeline.PutJobFailureResultInput, _ ...func(*codepipeline.Options)) (*codepipeline.PutJobFailureResultOutput, error) {
	return &codepipeline.PutJobFailureResultOutput{}, m.err
}

func newCodePipelineSecret(days int) *vault.Secret {
	return &vault.Secret{
		Path:      "secret/pipeline/token",
		ExpiresAt: time.Now().Add(time.Duration(days) * 24 * time.Hour),
	}
}

func TestNewCodePipelineNotifier_MissingJobID(t *testing.T) {
	_, err := NewCodePipelineNotifier(vaultconfig.CodePipelineConfig{Region: "us-east-1"})
	if err == nil {
		t.Fatal("expected error for missing job_id")
	}
}

func TestNewCodePipelineNotifier_MissingRegion(t *testing.T) {
	_, err := NewCodePipelineNotifier(vaultconfig.CodePipelineConfig{JobID: "abc-123"})
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewCodePipelineNotifier_Valid(t *testing.T) {
	n := newCodePipelineNotifierWithClient(&mockCodePipelineClient{}, "job-id-123")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestCodePipelineNotifier_Notify_ExpiringSoon(t *testing.T) {
	n := newCodePipelineNotifierWithClient(&mockCodePipelineClient{}, "job-id-123")
	if err := n.Notify(context.Background(), newCodePipelineSecret(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCodePipelineNotifier_Notify_Expired(t *testing.T) {
	n := newCodePipelineNotifierWithClient(&mockCodePipelineClient{}, "job-id-123")
	if err := n.Notify(context.Background(), newCodePipelineSecret(-1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCodePipelineNotifier_Notify_ClientError(t *testing.T) {
	n := newCodePipelineNotifierWithClient(&mockCodePipelineClient{err: errors.New("aws error")}, "job-id-123")
	err := n.Notify(context.Background(), newCodePipelineSecret(3))
	if err == nil {
		t.Fatal("expected error from client")
	}
}
