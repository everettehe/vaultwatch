package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type secretsManagerClient interface {
	PutSecretValue(ctx context.Context, params *secretsmanager.PutSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.PutSecretValueOutput, error)
}

// SecretsManagerNotifier writes expiration events to AWS Secrets Manager as secret metadata.
type SecretsManagerNotifier struct {
	client    secretsManagerClient
	secretID  string
	region    string
}

// NewSecretsManagerNotifier creates a SecretsManagerNotifier using the default AWS config.
func NewSecretsManagerNotifier(secretID, region string) (*SecretsManagerNotifier, error) {
	if secretID == "" {
		return nil, fmt.Errorf("awssecretsmanager: secret_id is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("awssecretsmanager: failed to load AWS config: %w", err)
	}
	client := secretsmanager.NewFromConfig(cfg)
	return newSecretsManagerNotifierWithClient(client, secretID, region), nil
}

func newSecretsManagerNotifierWithClient(client secretsManagerClient, secretID, region string) *SecretsManagerNotifier {
	return &SecretsManagerNotifier{
		client:   client,
		secretID: secretID,
		region:   region,
	}
}

// Notify sends a secret expiration event to AWS Secrets Manager.
func (n *SecretsManagerNotifier) Notify(ctx context.Context, secret vault.Secret) error {
	msg := FormatMessage(secret)
	payload := map[string]string{
		"path":    secret.Path,
		"status":  msg.Status,
		"subject": msg.Subject,
		"body":    msg.Body,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("awssecretsmanager: failed to marshal payload: %w", err)
	}
	_, err = n.client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(n.secretID),
		SecretString: aws.String(string(data)),
	})
	if err != nil {
		return fmt.Errorf("awssecretsmanager: failed to put secret value: %w", err)
	}
	return nil
}
