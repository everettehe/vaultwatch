package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// snsSMSClient defines the subset of the SNS API used by SNSSMSNotifier.
type snsSMSClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSSMSNotifier sends SMS alerts via AWS SNS direct phone number publishing.
type SNSSMSNotifier struct {
	client      snsSMSClient
	phoneNumber string
	senderID    string
}

// NewSNSSMSNotifier creates a new SNSSMSNotifier using the default AWS config.
func NewSNSSMSNotifier(phoneNumber, senderID string) (*SNSSMSNotifier, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("sns_sms: phone_number is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sns_sms: failed to load AWS config: %w", err)
	}
	client := sns.NewFromConfig(cfg)
	return newSNSSMSNotifierWithClient(client, phoneNumber, senderID)
}

func newSNSSMSNotifierWithClient(client snsSMSClient, phoneNumber, senderID string) (*SNSSMSNotifier, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("sns_sms: phone_number is required")
	}
	return &SNSSMSNotifier{
		client:      client,
		phoneNumber: phoneNumber,
		senderID:    senderID,
	}, nil
}

// Notify sends an SMS message for the given secret via AWS SNS.
func (n *SNSSMSNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg, _ := FormatMessage(secret)
	input := &sns.PublishInput{
		Message:     aws.String(msg.Body),
		PhoneNumber: aws.String(n.phoneNumber),
		MessageAttributes: map[string]snstypes.MessageAttributeValue{},
	}
	if n.senderID != "" {
		input.MessageAttributes["AWS.SNS.SMS.SenderID"] = snstypes.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(n.senderID),
		}
	}
	_, err := n.client.Publish(ctx, input)
	if err != nil {
		return fmt.Errorf("sns_sms: failed to send SMS: %w", err)
	}
	return nil
}
