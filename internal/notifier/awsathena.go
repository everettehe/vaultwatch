package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"

	vaultwatch "github.com/yourusername/vaultwatch/internal/vault"
)

// athenaClient defines the subset of the Athena API used by this notifier.
type athenaClient interface {
	StartQueryExecution(ctx context.Context, params *athena.StartQueryExecutionInput, optFns ...func(*athena.Options)) (*athena.StartQueryExecutionOutput, error)
}

// AthenaNotifier submits a query to Amazon Athena when a secret is expiring.
type AthenaNotifier struct {
	client     athenaClient
	database   string
	workgroup  string
	outputLoc  string
}

// NewAthenaNotifier creates an AthenaNotifier using ambient AWS credentials.
func NewAthenaNotifier(database, workgroup, outputLoc, region string) (*AthenaNotifier, error) {
	if database == "" {
		return nil, fmt.Errorf("athena: database is required")
	}
	if outputLoc == "" {
		return nil, fmt.Errorf("athena: output_location is required")
	}
	if region == "" {
		return nil, fmt.Errorf("athena: region is required")
	}
	if workgroup == "" {
		workgroup = "primary"
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("athena: failed to load AWS config: %w", err)
	}
	return newAthenaNotifierWithClient(athena.NewFromConfig(cfg), database, workgroup, outputLoc), nil
}

func newAthenaNotifierWithClient(client athenaClient, database, workgroup, outputLoc string) *AthenaNotifier {
	return &AthenaNotifier{
		client:    client,
		database:  database,
		workgroup: workgroup,
		outputLoc: outputLoc,
	}
}

// Notify submits an Athena query recording the secret expiration event.
func (n *AthenaNotifier) Notify(ctx context.Context, secret *vaultwatch.Secret) error {
	msg, _ := FormatMessage(secret)
	query := fmt.Sprintf(
		"INSERT INTO vault_secret_events (path, message, event_time) VALUES ('%s', '%s', '%s')",
		secret.Path,
		msg.Body,
		time.Now().UTC().Format(time.RFC3339),
	)
	_, err := n.client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
		QueryString: aws.String(query),
		QueryExecutionContext: &types.QueryExecutionContext{
			Database: aws.String(n.database),
		},
		WorkGroup: aws.String(n.workgroup),
		ResultConfiguration: &types.ResultConfiguration{
			OutputLocation: aws.String(n.outputLoc),
		},
	})
	if err != nil {
		return fmt.Errorf("athena: failed to start query execution: %w", err)
	}
	return nil
}
