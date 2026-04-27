package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// AWSConfigClient defines the subset of AWS Config API used by this notifier.
type AWSConfigClient interface {
	PutEvaluations(ctx context.Context, params *configservice.PutEvaluationsInput, optFns ...func(*configservice.Options)) (*configservice.PutEvaluationsOutput, error)
}

// AWSConfigNotifier sends compliance evaluations to AWS Config.
type AWSConfigNotifier struct {
	client    AWSConfigClient
	resultToken string
	region    string
}

// NewAWSConfigNotifier creates a new AWSConfigNotifier using default AWS credentials.
func NewAWSConfigNotifier(resultToken, region string) (*AWSConfigNotifier, error) {
	if resultToken == "" {
		return nil, fmt.Errorf("awsconfig: result token is required")
	}
	if region == "" {
		return nil, fmt.Errorf("awsconfig: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("awsconfig: failed to load AWS config: %w", err)
	}
	return newAWSConfigNotifierWithClient(configservice.NewFromConfig(cfg), resultToken, region)
}

func newAWSConfigNotifierWithClient(client AWSConfigClient, resultToken, region string) (*AWSConfigNotifier, error) {
	if resultToken == "" {
		return nil, fmt.Errorf("awsconfig: result token is required")
	}
	if region == "" {
		return nil, fmt.Errorf("awsconfig: region is required")
	}
	return &AWSConfigNotifier{client: client, resultToken: resultToken, region: region}, nil
}

// Notify sends a compliance evaluation to AWS Config based on secret expiration status.
func (n *AWSConfigNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	compliance := types.ComplianceTypeCompliant
	if secret.IsExpired() || secret.IsExpiringSoon(0) {
		compliance = types.ComplianceTypeNonCompliant
	}

	_, err := n.client.PutEvaluations(ctx, &configservice.PutEvaluationsInput{
		ResultToken: aws.String(n.resultToken),
		Evaluations: []types.Evaluation{
			{
				ComplianceResourceId:   aws.String(secret.Path),
				ComplianceResourceType: aws.String("AWS::::Account"),
				ComplianceType:         compliance,
				OrderingTimestamp:      aws.Time(secret.ExpiresAt),
				Annotation:             aws.String(FormatMessage(secret).Body),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("awsconfig: failed to put evaluation: %w", err)
	}
	return nil
}
