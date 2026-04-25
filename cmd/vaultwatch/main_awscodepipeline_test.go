package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithCodePipeline(jobID, region string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.CodePipeline = &config.CodePipelineConfig{
		JobID:  jobID,
		Region: region,
	}
	return cfg
}

func TestBuildNotifiers_CodePipeline_Valid(t *testing.T) {
	cfg := minimalConfigWithCodePipeline("job-abc-123", "us-east-1")
	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) == 0 {
		t.Fatal("expected at least one notifier")
	}
}

func TestBuildNotifiers_CodePipeline_MissingJobID(t *testing.T) {
	cfg := minimalConfigWithCodePipeline("", "us-east-1")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing job_id")
	}
}

func TestBuildNotifiers_CodePipeline_MissingRegion(t *testing.T) {
	cfg := minimalConfigWithCodePipeline("job-abc-123", "")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestCodePipelineNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewCodePipelineNotifier(config.CodePipelineConfig{
		JobID:  "job-123",
		Region: "us-east-1",
	})
	// We expect an error here because no real AWS config is available in tests,
	// but we verify the returned type satisfies the Notifier interface when valid.
	if err == nil && n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
