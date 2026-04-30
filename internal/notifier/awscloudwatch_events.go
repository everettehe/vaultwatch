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
)

type cloudWatchEventsClient interface {
	PutEvents(ctx context.Context, params *cloudwatchevents.PutEventsInput, optFns ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error)
}

// CloudWatchEventsNotifier sends vault secret expiration events to AWS CloudWatch Events.
type CloudWatchEventsNotifier struct {
	client     cloudWatchEventsClient
	eventBus   string
	source     string
	detailType string
	region     string
}

// NewCloudWatchEventsNotifier creates a new CloudWatchEventsNotifier.
func NewCloudWatchEventsNotifier(eventBus, source, detailType, region string) (*CloudWatchEventsNotifier, error) {
	if eventBus == "" {
		return nil, fmt.Errorf("cloudwatchevents: event bus name is required")
	}
	if source == "" {
		source = "vaultwatch"
	}
	if detailType == "" {
		detailType = "SecretExpiration"
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchevents: failed to load AWS config: %w", err)
	}
	client := cloudwatchevents.NewFromConfig(cfg)
	return newCloudWatchEventsNotifierWithClient(client, eventBus, source, detailType, region)
}

func newCloudWatchEventsNotifierWithClient(client cloudWatchEventsClient, eventBus, source, detailType, region string) (*CloudWatchEventsNotifier, error) {
	if eventBus == "" {
		return nil, fmt.Errorf("cloudwatchevents: event bus name is required")
	}
	if source == "" {
		source = "vaultwatch"
	}
	if detailType == "" {
		detailType = "SecretExpiration"
	}
	return &CloudWatchEventsNotifier{
		client:     client,
		eventBus:   eventBus,
		source:     source,
		detailType: detailType,
		region:     region,
	}, nil
}

// Notify sends a CloudWatch event for the given secret.
func (n *CloudWatchEventsNotifier) Notify(ctx context.Context, secret Secret) error {
	msg, _ := FormatMessage(secret)
	detail := map[string]interface{}{
		"path":       secret.Path,
		"days_until": secret.DaysUntilExpiration(),
		"expired":    secret.IsExpired(),
		"message":    msg.Body,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}
	detailJSON, err := json.Marshal(detail)
	if err != nil {
		return fmt.Errorf("cloudwatchevents: failed to marshal detail: %w", err)
	}
	_, err = n.client.PutEvents(ctx, &cloudwatchevents.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				EventBusName: aws.String(n.eventBus),
				Source:       aws.String(n.source),
				DetailType:   aws.String(n.detailType),
				Detail:       aws.String(string(detailJSON)),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("cloudwatchevents: failed to put event: %w", err)
	}
	return nil
}
