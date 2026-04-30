package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// PinpointSender abstracts the AWS Pinpoint SendMessages API.
type PinpointSender interface {
	SendMessages(ctx context.Context, params *pinpoint.SendMessagesInput, optFns ...func(*pinpoint.Options)) (*pinpoint.SendMessagesOutput, error)
}

// PinpointNotifier sends vault secret expiration alerts via AWS Pinpoint SMS.
type PinpointNotifier struct {
	client    PinpointSender
	appID     string
	origNumber string
	destNumber string
	region    string
}

// NewPinpointNotifier creates a PinpointNotifier using the provided config values.
func NewPinpointNotifier(appID, origNumber, destNumber, region string) (*PinpointNotifier, error) {
	if appID == "" {
		return nil, fmt.Errorf("pinpoint: app_id is required")
	}
	if destNumber == "" {
		return nil, fmt.Errorf("pinpoint: dest_number is required")
	}
	if region == "" {
		return nil, fmt.Errorf("pinpoint: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("pinpoint: failed to load AWS config: %w", err)
	}
	return newPinpointNotifierWithClient(pinpoint.NewFromConfig(cfg), appID, origNumber, destNumber, region)
}

func newPinpointNotifierWithClient(client PinpointSender, appID, origNumber, destNumber, region string) (*PinpointNotifier, error) {
	if appID == "" {
		return nil, fmt.Errorf("pinpoint: app_id is required")
	}
	if destNumber == "" {
		return nil, fmt.Errorf("pinpoint: dest_number is required")
	}
	if region == "" {
		return nil, fmt.Errorf("pinpoint: region is required")
	}
	return &PinpointNotifier{
		client:     client,
		appID:      appID,
		origNumber: origNumber,
		destNumber: destNumber,
		region:     region,
	}, nil
}

// Notify sends an SMS message via AWS Pinpoint for the given secret.
func (n *PinpointNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg := FormatMessage(secret)
	input := &pinpoint.SendMessagesInput{
		ApplicationId: aws.String(n.appID),
		MessageRequest: &types.MessageRequest{
			Addresses: map[string]types.AddressConfiguration{
				n.destNumber: {ChannelType: types.ChannelTypeSms},
			},
			MessageConfiguration: &types.DirectMessageConfiguration{
				SMSMessage: &types.SMSMessage{
					Body:             aws.String(msg.Body),
					OriginationNumber: aws.String(n.origNumber),
					MessageType:      types.MessageTypeTransactional,
				},
			},
		},
	}
	_, err := n.client.SendMessages(ctx, input)
	if err != nil {
		return fmt.Errorf("pinpoint: failed to send message: %w", err)
	}
	return nil
}
