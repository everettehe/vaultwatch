package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// eventBridgeClient is the subset of the EventBridge API we use.
type eventBridgeClient interface {
	PutEvents(ctx context.Context, params *eventbridge.PutEventsInput, optFns ...func(*eventbridge.Options)) (*eventbridge.PutEventsOutput, error)
}

// EventBridgeNotifier sends vault secret expiry events to AWS EventBridge.
type EventBridgeNotifier struct {
	client    eventBridgeClient
	eventBus  string
	source    string
	detailType string
}

// NewEventBridgeNotifier creates an EventBridgeNotifier.
// eventBus may be "default" or a custom event bus ARN.
func NewEventBridgeNotifier(eventBus, source, detailType string) (*EventBridgeNotifier, error) {
	if eventBus == "" {
		return nil, fmt.Errorf("eventbridge: event bus name or ARN is required")
	}
	if source == "" {
		source = "vaultwatch"
	}
	if detailType == "" {
		detailType = "VaultSecretExpiry"
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("eventbridge: failed to load AWS config: %w", err)
	}
	return &EventBridgeNotifier{
		client:     eventbridge.NewFromConfig(cfg),
		eventBus:   eventBus,
		source:     source,
		detailType: detailType,
	}, nil
}

func newEventBridgeNotifierWithClient(client eventBridgeClient, eventBus, source, detailType string) *EventBridgeNotifier {
	return &EventBridgeNotifier{
		client:     client,
		eventBus:   eventBus,
		source:     source,
		detailType: detailType,
	}
}

// Notify sends a PutEvents call to AWS EventBridge.
func (n *EventBridgeNotifier) Notify(s *vault.Secret) error {
	msg := FormatMessage(s)
	detail := map[string]interface{}{
		"path":       s.Path,
		"expires_at": s.ExpiresAt.Format(time.RFC3339),
		"days_left":  s.DaysUntilExpiration(),
		"summary":    msg.Body,
	}
	detailJSON, err := json.Marshal(detail)
	if err != nil {
		return fmt.Errorf("eventbridge: failed to marshal detail: %w", err)
	}
	_, err = n.client.PutEvents(context.Background(), &eventbridge.PutEventsInput{
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
		return fmt.Errorf("eventbridge: put events failed: %w", err)
	}
	return nil
}
