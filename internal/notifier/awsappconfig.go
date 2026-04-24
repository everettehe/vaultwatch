package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appconfigdata"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// appConfigDataClient defines the interface used for AppConfig interactions.
type appConfigDataClient interface {
	StartConfigurationSession(ctx context.Context, params *appconfigdata.StartConfigurationSessionInput, optFns ...func(*appconfigdata.Options)) (*appconfigdata.StartConfigurationSessionOutput, error)
}

// AppConfigNotifier sends a notification by starting a new AppConfig
// configuration session, signalling that a secret rotation is required.
type AppConfigNotifier struct {
	client      appConfigDataClient
	application string
	environment string
	profile     string
	region      string
}

// NewAppConfigNotifier creates an AppConfigNotifier using real AWS credentials.
func NewAppConfigNotifier(application, environment, profile, region string) (*AppConfigNotifier, error) {
	if application == "" {
		return nil, fmt.Errorf("appconfig: application is required")
	}
	if environment == "" {
		return nil, fmt.Errorf("appconfig: environment is required")
	}
	if profile == "" {
		return nil, fmt.Errorf("appconfig: profile is required")
	}
	if region == "" {
		region = "us-east-1"
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("appconfig: failed to load AWS config: %w", err)
	}
	return newAppConfigNotifierWithClient(appconfigdata.NewFromConfig(cfg), application, environment, profile, region)
}

func newAppConfigNotifierWithClient(client appConfigDataClient, application, environment, profile, region string) (*AppConfigNotifier, error) {
	if application == "" {
		return nil, fmt.Errorf("appconfig: application is required")
	}
	if environment == "" {
		return nil, fmt.Errorf("appconfig: environment is required")
	}
	if profile == "" {
		return nil, fmt.Errorf("appconfig: profile is required")
	}
	return &AppConfigNotifier{
		client:      client,
		application: application,
		environment: environment,
		profile:     profile,
		region:      region,
	}, nil
}

// Notify starts a new AppConfig configuration session to signal secret expiry.
func (n *AppConfigNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg := FormatMessage(secret)
	_, err := n.client.StartConfigurationSession(ctx, &appconfigdata.StartConfigurationSessionInput{
		ApplicationIdentifier:          aws.String(n.application),
		EnvironmentIdentifier:          aws.String(n.environment),
		ConfigurationProfileIdentifier: aws.String(n.profile),
		RequiredMinimumPollIntervalInSeconds: aws.Int32(15),
	})
	if err != nil {
		return fmt.Errorf("appconfig: failed to start configuration session for %s: %w: %s", secret.Path, err, msg.Body)
	}
	return nil
}
