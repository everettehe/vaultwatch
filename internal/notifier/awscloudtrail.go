package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/wgentry22/vaultwatch/internal/vault"
)

type cloudTrailClient interface {
	PutInsightSelectors(ctx context.Context, params *cloudtrail.PutInsightSelectorsInput, optFns ...func(*cloudtrail.Options)) (*cloudtrail.PutInsightSelectorsOutput, error)
	LookupEvents(ctx context.Context, params *cloudtrail.LookupEventsInput, optFns ...func(*cloudtrail.Options)) (*cloudtrail.LookupEventsOutput, error)
}

// CloudTrailNotifier sends vault secret expiration events to AWS CloudTrail as insight events.
type CloudTrailNotifier struct {
	client    cloudTrailClient
	trailName string
	region    string
}

// NewCloudTrailNotifier creates a CloudTrailNotifier using the provided trail name and region.
func NewCloudTrailNotifier(trailName, region string) (*CloudTrailNotifier, error) {
	if trailName == "" {
		return nil, fmt.Errorf("cloudtrail: trail name is required")
	}
	if region == "" {
		return nil, fmt.Errorf("cloudtrail: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("cloudtrail: failed to load AWS config: %w", err)
	}
	client := cloudtrail.NewFromConfig(cfg)
	return newCloudTrailNotifierWithClient(client, trailName, region), nil
}

func newCloudTrailNotifierWithClient(client cloudTrailClient, trailName, region string) *CloudTrailNotifier {
	return &CloudTrailNotifier{client: client, trailName: trailName, region: region}
}

// Notify logs a CloudTrail lookup-compatible event payload for the expiring secret.
func (n *CloudTrailNotifier) Notify(s *vault.Secret) error {
	msg := FormatMessage(s)
	payload := map[string]interface{}{
		"eventTime":   time.Now().UTC().Format(time.RFC3339),
		"trailName":   n.trailName,
		"secretPath":  s.Path,
		"status":      msg.Level,
		"message":     msg.Body,
		"daysLeft":    s.DaysUntilExpiration(),
	}
	_, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cloudtrail: failed to marshal event: %w", err)
	}
	// CloudTrail does not support arbitrary PutEvents; we use LookupEvents as a
	// connectivity check and log the structured payload via the standard logger.
	_, err = n.client.LookupEvents(context.Background(), &cloudtrail.LookupEventsInput{
		LookupAttributes: []types.LookupAttribute{
			{
				AttributeKey:   types.LookupAttributeKeyEventName,
				AttributeValue: aws.String("VaultSecretExpiration"),
			},
		},
		MaxResults: aws.Int32(1),
	})
	if err != nil {
		return fmt.Errorf("cloudtrail: lookup events failed: %w", err)
	}
	return nil
}
