package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockCloudfrontClient struct {
	err error
}

func (m *mockCloudfrontClient) CreateInvalidation(_ context.Context, _ *cloudfront.CreateInvalidationInput, _ ...func(*cloudfront.Options)) (*cloudfront.CreateInvalidationOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &cloudfront.CreateInvalidationOutput{}, nil
}

func newCloudfrontSecret() *vault.Secret {
	return &vault.Secret{
		Path:      "secret/myapp/tls",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
}

func TestNewCloudFrontNotifier_MissingDistributionID(t *testing.T) {
	_, err := NewCloudFrontNotifier("", "us-east-1", nil)
	if err == nil {
		t.Fatal("expected error for missing distribution_id")
	}
}

func TestNewCloudFrontNotifier_MissingRegion(t *testing.T) {
	_, err := NewCloudFrontNotifier("EDFDVBD6EXAMPLE", "", nil)
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewCloudFrontNotifier_DefaultPaths(t *testing.T) {
	client := &mockCloudfrontClient{}
	n := newCloudFrontNotifierWithClient(client, "EDFDVBD6EXAMPLE", nil)
	if len(n.paths) != 1 || n.paths[0] != "/*" {
		t.Errorf("expected default path '/*', got %v", n.paths)
	}
}

func TestCloudFrontNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockCloudfrontClient{}
	n := newCloudFrontNotifierWithClient(client, "EDFDVBD6EXAMPLE", []string{"/assets/*"})
	if err := n.Notify(context.Background(), newCloudfrontSecret()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloudFrontNotifier_Notify_ClientError(t *testing.T) {
	client := &mockCloudfrontClient{err: errors.New("aws error")}
	n := newCloudFrontNotifierWithClient(client, "EDFDVBD6EXAMPLE", []string{"/*"})
	if err := n.Notify(context.Background(), newCloudfrontSecret()); err == nil {
		t.Fatal("expected error from client")
	}
}

func TestCloudFrontNotifier_ImplementsInterface(t *testing.T) {
	client := &mockCloudfrontClient{}
	n := newCloudFrontNotifierWithClient(client, "EDFDVBD6EXAMPLE", nil)
	var _ Notifier = n
}
