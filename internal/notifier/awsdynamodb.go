package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// dynamoDBClient defines the interface used for DynamoDB operations.
type dynamoDBClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// DynamoDBNotifier writes secret expiration events to an AWS DynamoDB table.
type DynamoDBNotifier struct {
	client    dynamoDBClient
	tableName string
}

// NewDynamoDBNotifier creates a DynamoDBNotifier using the default AWS config.
func NewDynamoDBNotifier(tableName string) (*DynamoDBNotifier, error) {
	if tableName == "" {
		return nil, fmt.Errorf("dynamodb: table name is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("dynamodb: failed to load AWS config: %w", err)
	}
	return &DynamoDBNotifier{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: tableName,
	}, nil
}

// newDynamoDBNotifierWithClient creates a DynamoDBNotifier with a custom client (for testing).
func newDynamoDBNotifierWithClient(client dynamoDBClient, tableName string) (*DynamoDBNotifier, error) {
	if tableName == "" {
		return nil, fmt.Errorf("dynamodb: table name is required")
	}
	return &DynamoDBNotifier{client: client, tableName: tableName}, nil
}

// Notify writes the secret expiration event as a DynamoDB item.
func (n *DynamoDBNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg := FormatMessage(secret)
	_, err := n.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(n.tableName),
		Item: map[string]types.AttributeValue{
			"SecretPath": &types.AttributeValueMemberS{Value: secret.Path},
			"Timestamp":  &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
			"Status":     &types.AttributeValueMemberS{Value: msg.Subject},
			"Body":       &types.AttributeValueMemberS{Value: msg.Body},
			"DaysLeft":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", secret.DaysUntilExpiration())},
		},
	})
	if err != nil {
		return fmt.Errorf("dynamodb: failed to put item: %w", err)
	}
	return nil
}
