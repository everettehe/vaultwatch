package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/aws/aws-sdk-go-v2/service/codebuild/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type codeBuildClient interface {
	StartBuild(ctx context.Context, params *codebuild.StartBuildInput, optFns ...func(*codebuild.Options)) (*codebuild.StartBuildOutput, error)
}

// CodeBuildNotifier triggers an AWS CodeBuild project when a secret is expiring.
type CodeBuildNotifier struct {
	client      codeBuildClient
	projectName string
	region      string
}

// NewCodeBuildNotifier creates a CodeBuildNotifier using the provided project name and region.
func NewCodeBuildNotifier(projectName, region string) (*CodeBuildNotifier, error) {
	if projectName == "" {
		return nil, fmt.Errorf("codebuild: project name is required")
	}
	if region == "" {
		return nil, fmt.Errorf("codebuild: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("codebuild: failed to load AWS config: %w", err)
	}
	return newCodeBuildNotifierWithClient(codebuild.NewFromConfig(cfg), projectName, region), nil
}

func newCodeBuildNotifierWithClient(client codeBuildClient, projectName, region string) *CodeBuildNotifier {
	return &CodeBuildNotifier{client: client, projectName: projectName, region: region}
}

// Notify triggers a CodeBuild project build with environment variables describing the expiring secret.
func (n *CodeBuildNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg := FormatMessage(secret)
	_, err := n.client.StartBuild(ctx, &codebuild.StartBuildInput{
		ProjectName: aws.String(n.projectName),
		EnvironmentVariablesOverride: []types.EnvironmentVariable{
			{Name: aws.String("VAULTWATCH_SECRET_PATH"), Value: aws.String(secret.Path), Type: types.EnvironmentVariableTypePlaintext},
			{Name: aws.String("VAULTWATCH_MESSAGE"), Value: aws.String(msg.Body), Type: types.EnvironmentVariableTypePlaintext},
			{Name: aws.String("VAULTWATCH_SEVERITY"), Value: aws.String(string(msg.Severity)), Type: types.EnvironmentVariableTypePlaintext},
		},
	})
	if err != nil {
		return fmt.Errorf("codebuild: failed to start build for project %s: %w", n.projectName, err)
	}
	return nil
}
