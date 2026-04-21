package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"
)

// glueClient defines the subset of the AWS Glue API used by GlueNotifier.
type glueClient interface {
	CreateJob(ctx context.Context, params *glue.CreateJobInput, optFns ...func(*glue.Options)) (*glue.CreateJobOutput, error)
	StartJobRun(ctx context.Context, params *glue.StartJobRunInput, optFns ...func(*glue.Options)) (*glue.StartJobRunOutput, error)
}

// GlueNotifier triggers an AWS Glue job run when a secret is expiring or expired.
type GlueNotifier struct {
	client  glueClient
	jobName string
	region  string
}

// NewGlueNotifier creates a GlueNotifier using the provided job name and AWS region.
func NewGlueNotifier(jobName, region string) (*GlueNotifier, error) {
	if jobName == "" {
		return nil, fmt.Errorf("glue: job name is required")
	}
	if region == "" {
		return nil, fmt.Errorf("glue: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("glue: failed to load AWS config: %w", err)
	}
	return newGlueNotifierWithClient(glue.NewFromConfig(cfg), jobName, region), nil
}

func newGlueNotifierWithClient(client glueClient, jobName, region string) *GlueNotifier {
	return &GlueNotifier{client: client, jobName: jobName, region: region}
}

// Notify triggers a Glue job run with secret metadata passed as job arguments.
func (n *GlueNotifier) Notify(ctx context.Context, secret Secret) error {
	msg, _ := FormatMessage(secret)
	args := map[string]string{
		"--secret_path":    secret.Path,
		"--message":        msg.Body,
		"--triggered_at":   time.Now().UTC().Format(time.RFC3339),
		"--days_remaining": fmt.Sprintf("%d", secret.DaysUntilExpiration()),
	}
	_, err := n.client.StartJobRun(ctx, &glue.StartJobRunInput{
		JobName:   aws.String(n.jobName),
		Arguments: args,
		Timeout:   aws.Int32(10),
		WorkerType: types.WorkerTypeG1X,
		NumberOfWorkers: aws.Int32(2),
	})
	if err != nil {
		return fmt.Errorf("glue: failed to start job run for %q: %w", n.jobName, err)
	}
	return nil
}
