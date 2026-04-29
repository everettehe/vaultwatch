package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// bedrockAgentClient defines the interface for invoking a Bedrock agent.
type bedrockAgentClient interface {
	InvokeAgent(ctx context.Context, params *bedrockagentruntime.InvokeAgentInput, optFns ...func(*bedrockagentruntime.Options)) (*bedrockagentruntime.InvokeAgentOutput, error)
}

// BedrockAgentNotifier sends secret expiry notifications to an AWS Bedrock Agent.
type BedrockAgentNotifier struct {
	client     bedrockAgentClient
	agentID    string
	aliasID    string
	sessionID  string
	region     string
}

// NewBedrockAgentNotifier creates a BedrockAgentNotifier using the default AWS config.
func NewBedrockAgentNotifier(agentID, aliasID, region string) (*BedrockAgentNotifier, error) {
	if agentID == "" {
		return nil, fmt.Errorf("bedrock agent: agentID is required")
	}
	if aliasID == "" {
		return nil, fmt.Errorf("bedrock agent: aliasID is required")
	}
	if region == "" {
		return nil, fmt.Errorf("bedrock agent: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("bedrock agent: load aws config: %w", err)
	}
	client := bedrockagentruntime.NewFromConfig(cfg)
	return newBedrockAgentNotifierWithClient(client, agentID, aliasID, region), nil
}

func newBedrockAgentNotifierWithClient(client bedrockAgentClient, agentID, aliasID, region string) *BedrockAgentNotifier {
	return &BedrockAgentNotifier{
		client:    client,
		agentID:   agentID,
		aliasID:   aliasID,
		sessionID: "vaultwatch-session",
		region:    region,
	}
}

// Notify sends a secret expiry event to the configured Bedrock Agent.
func (n *BedrockAgentNotifier) Notify(ctx context.Context, secret vault.Secret) error {
	msg, err := json.Marshal(map[string]interface{}{
		"path":        secret.Path,
		"days_left":   secret.DaysUntilExpiration(),
		"expired":     secret.IsExpired(),
		"expiry_time": secret.ExpiresAt.Format("2006-01-02T15:04:05Z"),
	})
	if err != nil {
		return fmt.Errorf("bedrock agent: marshal payload: %w", err)
	}
	_, err = n.client.InvokeAgent(ctx, &bedrockagentruntime.InvokeAgentInput{
		AgentId:         aws.String(n.agentID),
		AgentAliasId:    aws.String(n.aliasID),
		SessionId:       aws.String(n.sessionID),
		InputText:       aws.String(string(msg)),
		EnableTrace:     aws.Bool(false),
		EndSession:      aws.Bool(false),
		SessionState:    &types.SessionState{},
	})
	if err != nil {
		return fmt.Errorf("bedrock agent: invoke agent: %w", err)
	}
	return nil
}
