package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type ecsClient interface {
	RunTask(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error)
}

// ECSNotifier triggers an ECS task when a secret is expiring or expired.
type ECSNotifier struct {
	client      ecsClient
	cluster     string
	taskDef     string
	container   string
	launchType  string
	region      string
}

// NewECSNotifier creates an ECSNotifier using the provided cluster, task definition, and region.
func NewECSNotifier(cluster, taskDef, container, launchType, region string) (*ECSNotifier, error) {
	if cluster == "" {
		return nil, fmt.Errorf("ecs notifier: cluster is required")
	}
	if taskDef == "" {
		return nil, fmt.Errorf("ecs notifier: task_def is required")
	}
	if region == "" {
		return nil, fmt.Errorf("ecs notifier: region is required")
	}
	if launchType == "" {
		launchType = "FARGATE"
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("ecs notifier: failed to load aws config: %w", err)
	}
	return newECSNotifierWithClient(ecs.NewFromConfig(cfg), cluster, taskDef, container, launchType, region)
}

func newECSNotifierWithClient(client ecsClient, cluster, taskDef, container, launchType, region string) (*ECSNotifier, error) {
	return &ECSNotifier{
		client:     client,
		cluster:    cluster,
		taskDef:    taskDef,
		container:  container,
		launchType: launchType,
		region:     region,
	}, nil
}

// Notify triggers an ECS task run with the secret path and status as environment overrides.
func (n *ECSNotifier) Notify(ctx context.Context, s vault.Secret) error {
	msg, err := json.Marshal(map[string]string{
		"path":   s.Path,
		"status": FormatMessage(s).Subject,
	})
	if err != nil {
		return fmt.Errorf("ecs notifier: marshal payload: %w", err)
	}

	input := &ecs.RunTaskInput{
		Cluster:        aws.String(n.cluster),
		TaskDefinition: aws.String(n.taskDef),
		LaunchType:     types.LaunchType(n.launchType),
		Overrides: &types.TaskOverride{
			ContainerOverrides: []types.ContainerOverride{
				{
					Name: aws.String(n.container),
					Environment: []types.KeyValuePair{
						{Name: aws.String("VAULTWATCH_PAYLOAD"), Value: aws.String(string(msg))},
					},
				},
			},
		},
	}
	_, err = n.client.RunTask(ctx, input)
	if err != nil {
		return fmt.Errorf("ecs notifier: run task: %w", err)
	}
	return nil
}
