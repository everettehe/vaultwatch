package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type mockEKSClient struct {
	createAddonFn func(ctx context.Context, params *eks.CreateAddonInput, optFns ...func(*eks.Options)) (*eks.CreateAddonOutput, error)
}

func (m *mockEKSClient) CreateAddon(ctx context.Context, params *eks.CreateAddonInput, optFns ...func(*eks.Options)) (*eks.CreateAddonOutput, error) {
	return m.createAddonFn(ctx, params, optFns...)
}

func newEKSSecret() *vault.Secret {
	return &vault.Secret{
		Path:      "secret/eks-token",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
}

func TestNewEKSNotifier_MissingCluster(t *testing.T) {
	_, err := NewEKSNotifier("", "vpc-cni", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing cluster name")
	}
}

func TestNewEKSNotifier_MissingAddon(t *testing.T) {
	_, err := NewEKSNotifier("my-cluster", "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing addon name")
	}
}

func TestNewEKSNotifier_MissingRegion(t *testing.T) {
	_, err := NewEKSNotifier("my-cluster", "vpc-cni", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewEKSNotifier_Valid(t *testing.T) {
	client := &mockEKSClient{}
	n := newEKSNotifierWithClient(client, "my-cluster", "vpc-cni", "us-east-1")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	if n.clusterName != "my-cluster" {
		t.Errorf("expected clusterName %q, got %q", "my-cluster", n.clusterName)
	}
	if n.addonName != "vpc-cni" {
		t.Errorf("expected addonName %q, got %q", "vpc-cni", n.addonName)
	}
}

func TestEKSNotifier_Notify_ExpiringSoon(t *testing.T) {
	var capturedInput *eks.CreateAddonInput
	client := &mockEKSClient{
		createAddonFn: func(ctx context.Context, params *eks.CreateAddonInput, optFns ...func(*eks.Options)) (*eks.CreateAddonOutput, error) {
			capturedInput = params
			return &eks.CreateAddonOutput{}, nil
		},
	}

	n := newEKSNotifierWithClient(client, "prod-cluster", "vpc-cni", "us-west-2")
	secret := newEKSSecret()

	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedInput == nil {
		t.Fatal("expected CreateAddon to be called")
	}
	if capturedInput.Tags["vaultwatch:secret-path"] != secret.Path {
		t.Errorf("expected secret path tag %q, got %q", secret.Path, capturedInput.Tags["vaultwatch:secret-path"])
	}
}

func TestEKSNotifier_Notify_Error(t *testing.T) {
	client := &mockEKSClient{
		createAddonFn: func(ctx context.Context, params *eks.CreateAddonInput, optFns ...func(*eks.Options)) (*eks.CreateAddonOutput, error) {
			return nil, errors.New("access denied")
		},
	}

	n := newEKSNotifierWithClient(client, "prod-cluster", "vpc-cni", "us-west-2")
	err := n.Notify(context.Background(), newEKSSecret())
	if err == nil {
		t.Fatal("expected error from failed CreateAddon call")
	}
}
