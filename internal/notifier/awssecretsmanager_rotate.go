package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	vaultsecret "github.com/youorg/vaultwatch/internal/vault"
)

// rotateClient defines the interface for triggering AWS Secrets Manager rotation.
type rotateClient interface {
	RotateSecret(ctx context.Context, params *secretsmanager.RotateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.RotateSecretOutput, error)
}

// SecretsManagerRotateNotifier triggers an AWS Secrets Manager rotation when a
// secret is expiring or expired.
type SecretsManagerRotateNotifier struct {
	client   rotateClient
	secretID string
}

// NewSecretsManagerRotateNotifier creates a new notifier that triggers rotation
// via AWS Secrets Manager. secretID is the ARN or name of the managed secret.
func NewSecretsManagerRotateNotifier(secretID, region string) (*SecretsManagerRotateNotifier, error) {
	if secretID == "" {
		return nil, fmt.Errorf("secretsmanager rotate notifier: secret_id is required")
	}
	if region == "" {
		return nil, fmt.Errorf("secretsmanager rotate notifier: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("secretsmanager rotate notifier: failed to load AWS config: %w", err)
	}
	return newSecretsManagerRotateNotifierWithClient(
		secretID,
		secretsmgr.NewFromConfig(cfg),
	), nil
}

func newSecretsManagerRotateNotifierWithClient(secretID string, c rotateClient) *SecretsManagerRotateNotifier {
	return &SecretsManagerRotateNotifier{client: c, secretID: secretID}
}

// Notify triggers an AWS Secrets Manager rotation for the configured secret.
func (n *SecretsManagerRotateNotifier) Notify(ctx context.Context, s *vaultsecret.Secret) error {
	_, err := n.client.RotateSecret(ctx, &secretsmanager.RotateSecretInput{
		SecretId: aws.String(n.secretID),
		ClientRequestToken: aws.String(fmt.Sprintf("vaultwatch-%s", s.Path)),
	})
	if err != nil {
		return fmt.Errorf("secretsmanager rotate notifier: failed to rotate secret %q: %w", n.secretID, err)
	}
	return nil
}
