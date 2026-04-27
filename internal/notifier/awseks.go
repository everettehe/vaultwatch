package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type eksClient interface {
	CreateAddon(ctx context.Context, params *eks.CreateAddonInput, optFns ...func(*eks.Options)) (*eks.CreateAddonOutput, error)
}

// EKSNotifier triggers an EKS addon update as a notification mechanism,
// annotating the addon description with the expiring secret path.
type EKSNotifier struct {
	client      eksClient
	clusterName string
	addonName   string
	region      string
}

// NewEKSNotifier creates an EKSNotifier using the provided cluster and addon names.
func NewEKSNotifier(clusterName, addonName, region string) (*EKSNotifier, error) {
	if clusterName == "" {
		return nil, fmt.Errorf("eks notifier: cluster name is required")
	}
	if addonName == "" {
		return nil, fmt.Errorf("eks notifier: addon name is required")
	}
	if region == "" {
		return nil, fmt.Errorf("eks notifier: region is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("eks notifier: failed to load AWS config: %w", err)
	}

	return newEKSNotifierWithClient(eks.NewFromConfig(cfg), clusterName, addonName, region), nil
}

func newEKSNotifierWithClient(client eksClient, clusterName, addonName, region string) *EKSNotifier {
	return &EKSNotifier{
		client:      client,
		clusterName: clusterName,
		addonName:   addonName,
		region:      region,
	}
}

// Notify sends a notification by tagging the EKS addon with expiry metadata.
func (n *EKSNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg := FormatMessage(secret)

	_, err := n.client.CreateAddon(ctx, &eks.CreateAddonInput{
		ClusterName: aws.String(n.clusterName),
		AddonName:   aws.String(n.addonName),
		Tags: map[string]string{
			"vaultwatch:secret-path":    secret.Path,
			"vaultwatch:alert-message": msg.Subject,
			"vaultwatch:status":        string(types.AddonStatusActive),
		},
	})
	if err != nil {
		return fmt.Errorf("eks notifier: failed to annotate addon %q on cluster %q: %w",
			n.addonName, n.clusterName, err)
	}

	return nil
}
