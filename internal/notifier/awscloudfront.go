package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type cloudfrontClient interface {
	CreateInvalidation(ctx context.Context, params *cloudfront.CreateInvalidationInput, optFns ...func(*cloudfront.Options)) (*cloudfront.CreateInvalidationOutput, error)
}

// CloudFrontNotifier triggers a CloudFront invalidation when a secret is expiring.
type CloudFrontNotifier struct {
	client       cloudfrontClient
	distributionID string
	paths        []string
}

// NewCloudFrontNotifier creates a CloudFrontNotifier using the default AWS config.
func NewCloudFrontNotifier(distributionID, region string, paths []string) (*CloudFrontNotifier, error) {
	if distributionID == "" {
		return nil, fmt.Errorf("cloudfront: distribution_id is required")
	}
	if region == "" {
		return nil, fmt.Errorf("cloudfront: region is required")
	}
	if len(paths) == 0 {
		paths = []string{"/*"}
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("cloudfront: failed to load AWS config: %w", err)
	}
	return newCloudFrontNotifierWithClient(cloudfront.NewFromConfig(cfg), distributionID, paths), nil
}

func newCloudFrontNotifierWithClient(client cloudfrontClient, distributionID string, paths []string) *CloudFrontNotifier {
	return &CloudFrontNotifier{
		client:         client,
		distributionID: distributionID,
		paths:          paths,
	}
}

// Notify creates a CloudFront invalidation referencing the expiring secret path.
func (n *CloudFrontNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	callerRef := fmt.Sprintf("vaultwatch-%d", time.Now().UnixNano())
	items := make([]string, len(n.paths))
	copy(items, n.paths)
	_, err := n.client.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(n.distributionID),
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: aws.String(callerRef),
			Paths: &types.Paths{
				Quantity: aws.Int32(int32(len(items))),
				Items:    items,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("cloudfront: create invalidation failed: %w", err)
	}
	return nil
}
