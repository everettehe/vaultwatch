package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// connectClient is the subset of the Connect API used by AWSConnectNotifier.
type connectClient interface {
	CreateContact(ctx context.Context, params *connect.CreateContactInput, optFns ...func(*connect.Options)) (*connect.CreateContactOutput, error)
}

// AWSConnectNotifier sends vault secret expiry alerts via Amazon Connect.
type AWSConnectNotifier struct {
	client     connectClient
	instanceID string
	contactFlow string
	queueID    string
}

// NewAWSConnectNotifier creates a new AWSConnectNotifier.
func NewAWSConnectNotifier(instanceID, contactFlow, queueID, region string) (*AWSConnectNotifier, error) {
	if instanceID == "" {
		return nil, fmt.Errorf("awsconnect: instance ID is required")
	}
	if contactFlow == "" {
		return nil, fmt.Errorf("awsconnect: contact flow ID is required")
	}
	if region == "" {
		return nil, fmt.Errorf("awsconnect: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("awsconnect: failed to load AWS config: %w", err)
	}
	return newAWSConnectNotifierWithClient(connect.NewFromConfig(cfg), instanceID, contactFlow, queueID)
}

func newAWSConnectNotifierWithClient(client connectClient, instanceID, contactFlow, queueID string) (*AWSConnectNotifier, error) {
	return &AWSConnectNotifier{
		client:      client,
		instanceID:  instanceID,
		contactFlow: contactFlow,
		queueID:     queueID,
	}, nil
}

// Notify sends an alert about a secret's expiry status via Amazon Connect.
func (n *AWSConnectNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg := FormatMessage(secret)
	_, err := n.client.CreateContact(ctx, &connect.CreateContactInput{
		InstanceId:            aws.String(n.instanceID),
		ContactFlowId:         aws.String(n.contactFlow),
		QueueId:               aws.String(n.queueID),
		Description:           aws.String(msg.Body),
		Name:                  aws.String(msg.Subject),
	})
	if err != nil {
		return fmt.Errorf("awsconnect: failed to create contact: %w", err)
	}
	return nil
}
