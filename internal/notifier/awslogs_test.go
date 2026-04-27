package notifier

import (
	"context"
	"testing"
	"time"

	vaultconfig "github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func newAWSLogsSecret() vault.Secret {
	return vault.Secret{
		Path:      "secret/logs-test",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
}

func TestNewCloudWatchLogsV2Notifier_MissingLogGroup(t *testing.T) {
	_, err := NewCloudWatchLogsV2Notifier(vaultconfig.CloudWatchLogsConfig{})
	if err == nil {
		t.Fatal("expected error for missing log group")
	}
}

func TestNewCloudWatchLogsV2Notifier_DefaultLogStream(t *testing.T) {
	cfg := vaultconfig.CloudWatchLogsConfig{
		LogGroup: "vaultwatch-logs",
		Region:   "us-east-1",
	}
	n, err := NewCloudWatchLogsV2Notifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.logStream != "vaultwatch" {
		t.Errorf("expected default log stream 'vaultwatch', got %q", n.logStream)
	}
}

func TestNewCloudWatchLogsV2Notifier_CustomLogStream(t *testing.T) {
	cfg := vaultconfig.CloudWatchLogsConfig{
		LogGroup:  "vaultwatch-logs",
		LogStream: "custom-stream",
		Region:    "us-east-1",
	}
	n, err := NewCloudWatchLogsV2Notifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.logStream != "custom-stream" {
		t.Errorf("expected log stream 'custom-stream', got %q", n.logStream)
	}
}

func TestCloudWatchLogsV2Notifier_ImplementsInterface(t *testing.T) {
	cfg := vaultconfig.CloudWatchLogsConfig{
		LogGroup: "vaultwatch-logs",
		Region:   "us-east-1",
	}
	n, err := NewCloudWatchLogsV2Notifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ Notifier = n
}

func TestCloudWatchLogsV2Notifier_Notify_ReturnsErrorOnBadClient(t *testing.T) {
	n := &CloudWatchLogsV2Notifier{
		client:    nil,
		logGroup:  "vaultwatch-logs",
		logStream: "vaultwatch",
	}
	s := newAWSLogsSecret()
	err := n.Notify(context.Background(), s)
	if err == nil {
		t.Fatal("expected error when client is nil")
	}
}
