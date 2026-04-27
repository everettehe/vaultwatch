package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/wgentry22/vaultwatch/internal/vault"
)

type mockCloudTrailClient struct {
	err error
}

func (m *mockCloudTrailClient) PutInsightSelectors(_ context.Context, _ *cloudtrail.PutInsightSelectorsInput, _ ...func(*cloudtrail.Options)) (*cloudtrail.PutInsightSelectorsOutput, error) {
	return &cloudtrail.PutInsightSelectorsOutput{}, m.err
}

func (m *mockCloudTrailClient) LookupEvents(_ context.Context, _ *cloudtrail.LookupEventsInput, _ ...func(*cloudtrail.Options)) (*cloudtrail.LookupEventsOutput, error) {
	return &cloudtrail.LookupEventsOutput{}, m.err
}

func newCloudTrailSecret(daysLeft int) *vault.Secret {
	expiry := time.Now().Add(time.Duration(daysLeft) * 24 * time.Hour)
	return &vault.Secret{Path: "secret/cloudtrail", Expiration: expiry}
}

func TestNewCloudTrailNotifier_MissingTrailName(t *testing.T) {
	_, err := NewCloudTrailNotifier("", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing trail name")
	}
}

func TestNewCloudTrailNotifier_MissingRegion(t *testing.T) {
	_, err := NewCloudTrailNotifier("my-trail", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewCloudTrailNotifier_Valid(t *testing.T) {
	n := newCloudTrailNotifierWithClient(&mockCloudTrailClient{}, "my-trail", "us-east-1")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	if n.trailName != "my-trail" {
		t.Errorf("expected trail name 'my-trail', got %q", n.trailName)
	}
	if n.region != "us-east-1" {
		t.Errorf("expected region 'us-east-1', got %q", n.region)
	}
}

func TestCloudTrailNotifier_Notify_ExpiringSoon(t *testing.T) {
	n := newCloudTrailNotifierWithClient(&mockCloudTrailClient{}, "my-trail", "us-east-1")
	s := newCloudTrailSecret(5)
	if err := n.Notify(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloudTrailNotifier_Notify_Expired(t *testing.T) {
	n := newCloudTrailNotifierWithClient(&mockCloudTrailClient{}, "my-trail", "us-east-1")
	s := newCloudTrailSecret(-1)
	if err := n.Notify(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloudTrailNotifier_Notify_ClientError(t *testing.T) {
	n := newCloudTrailNotifierWithClient(&mockCloudTrailClient{err: errors.New("api error")}, "my-trail", "us-east-1")
	s := newCloudTrailSecret(3)
	if err := n.Notify(s); err == nil {
		t.Fatal("expected error from client")
	}
}

func TestCloudTrailNotifier_ImplementsInterface(t *testing.T) {
	n := newCloudTrailNotifierWithClient(&mockCloudTrailClient{}, "my-trail", "us-east-1")
	var _ Notifier = n
}
