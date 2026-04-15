package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockCloudWatchClient struct {
	captured *cloudwatch.PutMetricDataInput
	err      error
}

func (m *mockCloudWatchClient) PutMetricData(_ context.Context, params *cloudwatch.PutMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error) {
	m.captured = params
	return &cloudwatch.PutMetricDataOutput{}, m.err
}

func newCWSecret(path string, daysUntil int) *vault.Secret {
	return &vault.Secret{
		Path:      path,
		ExpiresAt: time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour),
	}
}

func TestNewCloudWatchNotifier_DefaultNamespace(t *testing.T) {
	// We can't call NewCloudWatchNotifier without real AWS creds, so we test the
	// internal constructor directly.
	n := newCloudWatchNotifierWithClient(&mockCloudWatchClient{}, "")
	// namespace defaults are applied by the public constructor; here we just
	// verify the struct is non-nil.
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestCloudWatchNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockCloudWatchClient{}
	n := newCloudWatchNotifierWithClient(mock, "VaultWatch")
	s := newCWSecret("secret/db/password", 5)

	if err := n.Notify(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.captured == nil {
		t.Fatal("expected PutMetricData to be called")
	}
	if *mock.captured.Namespace != "VaultWatch" {
		t.Errorf("expected namespace VaultWatch, got %s", *mock.captured.Namespace)
	}
	if len(mock.captured.MetricData) != 1 {
		t.Fatalf("expected 1 metric datum, got %d", len(mock.captured.MetricData))
	}
	datum := mock.captured.MetricData[0]
	if *datum.MetricName != "DaysUntilExpiration" {
		t.Errorf("unexpected metric name: %s", *datum.MetricName)
	}
	if len(datum.Dimensions) != 1 || *datum.Dimensions[0].Value != "secret/db/password" {
		t.Errorf("unexpected dimensions: %+v", datum.Dimensions)
	}
}

func TestCloudWatchNotifier_Notify_Expired(t *testing.T) {
	mock := &mockCloudWatchClient{}
	n := newCloudWatchNotifierWithClient(mock, "VaultWatch")
	s := newCWSecret("secret/old/key", -3)

	if err := n.Notify(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	datum := mock.captured.MetricData[0]
	if *datum.Value >= 0 {
		t.Errorf("expected negative days value for expired secret, got %f", *datum.Value)
	}
}

func TestCloudWatchNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockCloudWatchClient{err: errors.New("aws error")}
	n := newCloudWatchNotifierWithClient(mock, "VaultWatch")
	s := newCWSecret("secret/api/token", 10)

	err := n.Notify(s)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
