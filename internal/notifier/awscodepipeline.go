package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	vaultconfig "github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type codePipelineClient interface {
	PutJobFailureResult(ctx context.Context, params *codepipeline.PutJobFailureResultInput, optFns ...func(*codepipeline.Options)) (*codepipeline.PutJobFailureResultOutput, error)
}

// CodePipelineNotifier sends a job failure result to AWS CodePipeline when a secret is expiring.
type CodePipelineNotifier struct {
	client codePipelineClient
	jobID  string
}

// NewCodePipelineNotifier creates a CodePipelineNotifier using the provided config.
func NewCodePipelineNotifier(cfg vaultconfig.CodePipelineConfig) (*CodePipelineNotifier, error) {
	if cfg.JobID == "" {
		return nil, fmt.Errorf("codepipeline: job_id is required")
	}
	if cfg.Region == "" {
		return nil, fmt.Errorf("codepipeline: region is required")
	}
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("codepipeline: failed to load AWS config: %w", err)
	}
	return newCodePipelineNotifierWithClient(codepipeline.NewFromConfig(awsCfg), cfg.JobID), nil
}

func newCodePipelineNotifierWithClient(client codePipelineClient, jobID string) *CodePipelineNotifier {
	return &CodePipelineNotifier{client: client, jobID: jobID}
}

// Notify sends a failure result to CodePipeline with the secret expiration message.
func (n *CodePipelineNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg := FormatMessage(secret)
	_, err := n.client.PutJobFailureResult(ctx, &codepipeline.PutJobFailureResultInput{
		JobId: aws.String(n.jobID),
		FailureDetails: &types.FailureDetails{
			Message: aws.String(fmt.Sprintf("%s: %s", msg.Subject, msg.Body)),
			Type:    types.FailureTypeJobFailed,
		},
	})
	if err != nil {
		return fmt.Errorf("codepipeline: failed to put job failure result: %w", err)
	}
	return nil
}
