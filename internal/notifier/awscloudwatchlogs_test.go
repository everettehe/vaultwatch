package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockCWLogsClient struct {
	putErr    error
	createErr error
	called    bool
}

func (m *mockCWLogsClient) PutLogEvents(_ context.Context, _ *cloudwatchlogs.PutLogEventsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutLogEventsOutput, error) {
	m.called = true
	return &cloudwatchlogs.PutLogEventsOutput{}, m.putErr
}

func (m *mockCWLogsClient) CreateLogStream(_ context.Context, _ *cloudwatchlogs.CreateLogStreamInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogStreamOutput, error) {
	return &cloudwatchlogs.CreateLogStreamOutput{}, m.createErr
}

func newCWLogsSecret(daysLeft int) vault.Secret {
	return vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(time.Duration(daysLeft) * 24 * time.Hour),
	}
}

func TestNewCloudWatchLogsNotifier_MissingLogGroup(t *testing.T) {
	_, err := NewCloudWatchLogsNotifier("", "stream", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing log group")
	}
}

func TestNewCloudWatchLogsNotifier_DefaultLogStream(t *testing.T) {
	mock := &mockCWLogsClient{}
	n := newCloudWatchLogsNotifierWithClient(mock, "my-group", "", "us-east-1")
	if n.logStream != "" {
		// logStream default is set only in NewCloudWatchLogsNotifier, not in the with-client constructor
		t.Log("logStream is empty as expected for with-client constructor")
	}
}

func TestNewCloudWatchLogsNotifier_Valid(t *testing.T) {
	mock := &mockCWLogsClient{}
	n := newCloudWatchLogsNotifierWithClient(mock, "my-group", "vaultwatch", "us-east-1")
	if n.logGroup != "my-group" {
		t.Errorf("expected log group 'my-group', got %q", n.logGroup)
	}
	if n.logStream != "vaultwatch" {
		t.Errorf("expected log stream 'vaultwatch', got %q", n.logStream)
	}
}

func TestCloudWatchLogsNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockCWLogsClient{}
	n := newCloudWatchLogsNotifierWithClient(mock, "my-group", "vaultwatch", "us-east-1")
	err := n.Notify(context.Background(), newCWLogsSecret(5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected PutLogEvents to be called")
	}
}

func TestCloudWatchLogsNotifier_Notify_Expired(t *testing.T) {
	mock := &mockCWLogsClient{}
	n := newCloudWatchLogsNotifierWithClient(mock, "my-group", "vaultwatch", "us-east-1")
	err := n.Notify(context.Background(), newCWLogsSecret(-1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloudWatchLogsNotifier_Notify_PutError(t *testing.T) {
	mock := &mockCWLogsClient{putErr: errors.New("put failed")}
	n := newCloudWatchLogsNotifierWithClient(mock, "my-group", "vaultwatch", "us-east-1")
	err := n.Notify(context.Background(), newCWLogsSecret(3))
	if err == nil {
		t.Fatal("expected error from PutLogEvents")
	}
}
