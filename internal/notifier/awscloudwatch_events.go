package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
	"github.com/younsl/vaultwatch/internal/vault"
)

type cloudWatchEventsClient interface {
	PutEvents(ctx context.Context, params *cloudwatchevents.PutEventsInput, optFns ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error)
}

// CloudWatchEventsNotifier sends Vault secret expiration events to Amazon EventBridge (CloudWatch Events).
type CloudWatchEventsNotifier struct {
	client    cloudWatchEventsClient
	eventBus  string
	source    string
	detailType string
}

// NewCloudWatchEventsNotifier creates a new CloudWatchEventsNotifier.
func NewCloudWatchEventsNotifier(eventBus, source, detailType, region string) (*CloudWatchEventsNotifier, error) {
	if eventBus == "" {
		return nil, fmt.Errorf("cloudwatch events: event bus name is required")
	}
	if region == "" {
		return nil, fmt.Errorf("cloudwatch events: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("cloudwatch events: failed to load AWS config: %w", err)
	}
	client := cloudwatchevents.NewFromConfig(cfg)
	return newCloudWatchEventsNotifierWithClient(client, eventBus, source, detailType)
}

func newCloudWatchEventsNotifierWithClient(client cloudWatchEventsClient, eventBus, source, detailType string) (*CloudWatchEventsNotifier, error) {
	if source == "" {
		source = "vaultwatch"
	}
	if detailType == "" {
		detailType = "VaultSecretExpiration"
	}
	return &CloudWatchEventsNotifier{
		client:     client,
		eventBus:   eventBus,
		source:     source,
		detailType: detailType,
	}, nil
}

// Notify sends the secret expiration event to CloudWatch Events.
func (n *CloudWatchEventsNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	payload := map[string]interface{}{
		"path":        secret.Path,
		"days_until":  secret.DaysUntilExpiration(),
		"expires_at":  secret.ExpiresAt.Format(time.RFC3339),
		"is_expired":  secret.IsExpired(),
	}
	detail, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cloudwatch events: failed to marshal event detail: %w", err)
	}
	_, err = n.client.PutEvents(ctx, &cloudwatchevents.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				EventBusName: aws.String(n.eventBus),
				Source:       aws.String(n.source),
				DetailType:   aws.String(n.detailType),
				Detail:       aws.String(string(detail)),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("cloudwatch events: failed to put event: %w", err)
	}
	return nil
}
