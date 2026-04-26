package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	vaultconfig "github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// cloudformationClient defines the subset of the CloudFormation API used by CloudFormationNotifier.
type cloudformationClient interface {
	CreateStack(ctx context.Context, params *cloudformation.CreateStackInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateStackOutput, error)
}

// CloudFormationNotifier triggers a CloudFormation stack creation as an alert
// mechanism when a Vault secret is expiring or has expired.
type CloudFormationNotifier struct {
	client    cloudformationClient
	stackName string
	templateURL string
	region    string
}

// NewCloudFormationNotifier creates a CloudFormationNotifier from the provided configuration.
// It returns an error if required fields are missing.
func NewCloudFormationNotifier(cfg *vaultconfig.CloudFormationConfig) (*CloudFormationNotifier, error) {
	if cfg.StackName == "" {
		return nil, fmt.Errorf("cloudformation notifier: stack_name is required")
	}
	if cfg.TemplateURL == "" {
		return nil, fmt.Errorf("cloudformation notifier: template_url is required")
	}
	if cfg.Region == "" {
		return nil, fmt.Errorf("cloudformation notifier: region is required")
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("cloudformation notifier: failed to load AWS config: %w", err)
	}

	client := cloudformation.NewFromConfig(awsCfg)
	return newCloudFormationNotifierWithClient(cfg, client)
}

func newCloudFormationNotifierWithClient(cfg *vaultconfig.CloudFormationConfig, client cloudformationClient) (*CloudFormationNotifier, error) {
	return &CloudFormationNotifier{
		client:      client,
		stackName:   cfg.StackName,
		templateURL: cfg.TemplateURL,
		region:      cfg.Region,
	}, nil
}

// Notify triggers a CloudFormation stack creation with parameters describing
// the expiring or expired Vault secret.
func (n *CloudFormationNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg := FormatMessage(secret)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Use a unique stack name suffix to avoid collisions on repeated alerts.
	stackName := fmt.Sprintf("%s-%d", n.stackName, time.Now().UnixMilli())

	input := &cloudformation.CreateStackInput{
		StackName:   aws.String(stackName),
		TemplateURL: aws.String(n.templateURL),
		Parameters: []types.Parameter{
			{
				ParameterKey:   aws.String("SecretPath"),
				ParameterValue: aws.String(secret.Path),
			},
			{
				ParameterKey:   aws.String("AlertMessage"),
				ParameterValue: aws.String(msg.Subject),
			},
			{
				ParameterKey:   aws.String("AlertTimestamp"),
				ParameterValue: aws.String(timestamp),
			},
		},
		Tags: []types.Tag{
			{
				Key:   aws.String("ManagedBy"),
				Value: aws.String("vaultwatch"),
			},
			{
				Key:   aws.String("SecretPath"),
				Value: aws.String(secret.Path),
			},
		},
	}

	_, err := n.client.CreateStack(ctx, input)
	if err != nil {
		return fmt.Errorf("cloudformation notifier: failed to create stack %q: %w", stackName, err)
	}

	return nil
}
