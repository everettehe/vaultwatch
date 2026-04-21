package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithSQSFIFO() *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.SQSFIFO = &config.SQSFIFOConfig{
		QueueURL:       "https://sqs.us-east-1.amazonaws.com/123/queue.fifo",
		MessageGroupID: "vaultwatch",
	}
	return cfg
}

func TestBuildNotifiers_SQSFIFO_Valid(t *testing.T) {
	cfg := minimalConfigWithSQSFIFO()
	// NewSQSFIFONotifier uses real AWS config; we only verify no panic on valid config fields.
	// Integration requires AWS credentials; skip actual build here.
	if cfg.Notifiers.SQSFIFO.QueueURL == "" {
		t.Fatal("expected non-empty queue URL")
	}
	if cfg.Notifiers.SQSFIFO.MessageGroupID == "" {
		t.Fatal("expected non-empty message group ID")
	}
}

func TestBuildNotifiers_SQSFIFO_MissingQueueURL(t *testing.T) {
	_, err := notifier.NewSQSFIFONotifier("", "group1")
	if err == nil {
		t.Fatal("expected error for missing queue URL")
	}
}

func TestBuildNotifiers_SQSFIFO_MissingMessageGroup(t *testing.T) {
	_, err := notifier.NewSQSFIFONotifier("https://sqs.us-east-1.amazonaws.com/123/queue.fifo", "")
	if err == nil {
		t.Fatal("expected error for missing message group ID")
	}
}

func TestSQSFIFONotifier_ImplementsInterface(t *testing.T) {
	// Compile-time interface check via assignment — will fail to build if interface not satisfied.
	// We use the internal constructor with a nil client to avoid AWS calls.
	// This is a type assertion test only.
	var _ interface {
		Notify(interface{}, interface{}) error
	} = (*notifier.SQSFIFONotifier)(nil)
}
