package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/glue"
)

type mockGlueClient struct {
	startJobRunFunc func(ctx context.Context, params *glue.StartJobRunInput, optFns ...func(*glue.Options)) (*glue.StartJobRunOutput, error)
}

func (m *mockGlueClient) CreateJob(ctx context.Context, params *glue.CreateJobInput, optFns ...func(*glue.Options)) (*glue.CreateJobOutput, error) {
	return &glue.CreateJobOutput{}, nil
}

func (m *mockGlueClient) StartJobRun(ctx context.Context, params *glue.StartJobRunInput, optFns ...func(*glue.Options)) (*glue.StartJobRunOutput, error) {
	if m.startJobRunFunc != nil {
		return m.startJobRunFunc(ctx, params, optFns...)
	}
	return &glue.StartJobRunOutput{}, nil
}

func newGlueSecret(daysUntil int) Secret {
	expiry := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return Secret{Path: "secret/glue/token", ExpiresAt: expiry}
}

func TestNewGlueNotifier_MissingJobName(t *testing.T) {
	_, err := NewGlueNotifier("", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing job name")
	}
}

func TestNewGlueNotifier_MissingRegion(t *testing.T) {
	_, err := NewGlueNotifier("my-job", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewGlueNotifier_Valid(t *testing.T) {
	client := &mockGlueClient{}
	n := newGlueNotifierWithClient(client, "rotation-job", "us-east-1")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	if n.jobName != "rotation-job" {
		t.Errorf("expected jobName %q, got %q", "rotation-job", n.jobName)
	}
}

func TestGlueNotifier_Notify_ExpiringSoon(t *testing.T) {
	var capturedInput *glue.StartJobRunInput
	client := &mockGlueClient{
		startJobRunFunc: func(ctx context.Context, params *glue.StartJobRunInput, _ ...func(*glue.Options)) (*glue.StartJobRunOutput, error) {
			capturedInput = params
			return &glue.StartJobRunOutput{}, nil
		},
	}
	n := newGlueNotifierWithClient(client, "rotation-job", "us-east-1")
	secret := newGlueSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedInput == nil {
		t.Fatal("expected StartJobRun to be called")
	}
	if *capturedInput.JobName != "rotation-job" {
		t.Errorf("expected job name %q, got %q", "rotation-job", *capturedInput.JobName)
	}
	if capturedInput.Arguments["--secret_path"] != "secret/glue/token" {
		t.Errorf("expected secret_path arg, got %q", capturedInput.Arguments["--secret_path"])
	}
}

func TestGlueNotifier_Notify_Expired(t *testing.T) {
	client := &mockGlueClient{}
	n := newGlueNotifierWithClient(client, "rotation-job", "us-east-1")
	secret := newGlueSecret(-1)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGlueNotifier_Notify_Error(t *testing.T) {
	client := &mockGlueClient{
		startJobRunFunc: func(_ context.Context, _ *glue.StartJobRunInput, _ ...func(*glue.Options)) (*glue.StartJobRunOutput, error) {
			return nil, errors.New("glue API error")
		},
	}
	n := newGlueNotifierWithClient(client, "rotation-job", "us-east-1")
	secret := newGlueSecret(3)
	if err := n.Notify(context.Background(), secret); err == nil {
		t.Fatal("expected error from Notify")
	}
}

func TestGlueNotifier_ImplementsInterface(t *testing.T) {
	client := &mockGlueClient{}
	n := newGlueNotifierWithClient(client, "job", "us-east-1")
	var _ Notifier = n
}
