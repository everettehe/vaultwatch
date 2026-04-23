package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/aws/aws-sdk-go-v2/service/emr/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type emrClient interface {
	AddJobFlowSteps(ctx context.Context, params *emr.AddJobFlowStepsInput, optFns ...func(*emr.Options)) (*emr.AddJobFlowStepsOutput, error)
}

// EMRNotifier sends vault secret expiration alerts as EMR job flow steps.
type EMRNotifier struct {
	client    emrClient
	clusterID string
	region    string
}

// NewEMRNotifier creates a new EMRNotifier using the provided cluster ID and region.
func NewEMRNotifier(clusterID, region string) (*EMRNotifier, error) {
	if clusterID == "" {
		return nil, fmt.Errorf("emr: cluster ID is required")
	}
	if region == "" {
		return nil, fmt.Errorf("emr: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("emr: failed to load AWS config: %w", err)
	}
	return &EMRNotifier{
		client:    emr.NewFromConfig(cfg),
		clusterID: clusterID,
		region:    region,
	}, nil
}

func newEMRNotifierWithClient(client emrClient, clusterID, region string) *EMRNotifier {
	return &EMRNotifier{client: client, clusterID: clusterID, region: region}
}

// Notify submits an EMR step that records the secret expiration event.
func (n *EMRNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg, err := json.Marshal(map[string]string{
		"path":       secret.Path,
		"expires_at": secret.ExpiresAt.Format("2006-01-02T15:04:05Z"),
		"message":    FormatMessage(secret),
	})
	if err != nil {
		return fmt.Errorf("emr: failed to marshal message: %w", err)
	}

	_, err = n.client.AddJobFlowSteps(ctx, &emr.AddJobFlowStepsInput{
		JobFlowId: aws.String(n.clusterID),
		Steps: []types.StepConfig{
			{
				Name:            aws.String(fmt.Sprintf("vaultwatch-alert-%s", secret.Path)),
				ActionOnFailure: types.ActionOnFailureContinue,
				HadoopJarStep: &types.HadoopJarStepConfig{
					Jar:  aws.String("command-runner.jar"),
					Args: []string{"echo", string(msg)},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("emr: failed to add job flow step: %w", err)
	}
	return nil
}
