package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	vaultwatch "github.com/yourusername/vaultwatch/internal/vault"
)

// ssmPutParameterAPI defines the interface used for putting SSM parameters.
type ssmPutParameterAPI interface {
	PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
}

// SSMNotifier writes secret expiration alerts as SSM Parameter Store entries.
type SSMNotifier struct {
	client    ssmPutParameterAPI
	paramPath string
}

// NewSSMNotifier creates an SSMNotifier using the default AWS configuration.
func NewSSMNotifier(paramPath string) (*SSMNotifier, error) {
	if paramPath == "" {
		return nil, fmt.Errorf("ssm: parameter path must not be empty")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ssm: failed to load AWS config: %w", err)
	}
	return &SSMNotifier{
		client:    ssm.NewFromConfig(cfg),
		paramPath: paramPath,
	}, nil
}

func newSSMNotifierWithClient(client ssmPutParameterAPI, paramPath string) (*SSMNotifier, error) {
	if paramPath == "" {
		return nil, fmt.Errorf("ssm: parameter path must not be empty")
	}
	return &SSMNotifier{client: client, paramPath: paramPath}, nil
}

// Notify writes the secret expiration status as a String parameter in SSM.
func (n *SSMNotifier) Notify(ctx context.Context, secret *vaultwatch.Secret) error {
	msg, _ := FormatMessage(secret)
	value := fmt.Sprintf("{\"path\":%q,\"days_until_expiration\":%d,\"message\":%q,\"timestamp\":%q}",
		secret.Path,
		int(secret.ExpiresAt.Sub(time.Now().UTC()).Hours()/24),
		msg.Body,
		time.Now().UTC().Format(time.RFC3339),
	)
	_, err := n.client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      aws.String(fmt.Sprintf("%s/%s", n.paramPath, sanitizeLabelValue(secret.Path))),
		Value:     aws.String(value),
		Type:      types.ParameterTypeString,
		Overwrite: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("ssm: failed to put parameter: %w", err)
	}
	return nil
}
