package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/yourusername/vaultwatch/internal/vault"
	"google.golang.org/api/option"
)

// GooglePubSubNotifier publishes secret expiration events to a Google Cloud Pub/Sub topic.
type GooglePubSubNotifier struct {
	projectID string
	topicID   string
	client    *pubsub.Client
}

type pubSubMessage struct {
	Path      string `json:"path"`
	Status    string `json:"status"`
	DaysLeft  int    `json:"days_until_expiration"`
	Message   string `json:"message"`
}

// NewGooglePubSubNotifier creates a GooglePubSubNotifier using Application Default Credentials.
func NewGooglePubSubNotifier(projectID, topicID string) (*GooglePubSubNotifier, error) {
	if projectID == "" {
		return nil, fmt.Errorf("googlepubsub: project_id is required")
	}
	if topicID == "" {
		return nil, fmt.Errorf("googlepubsub: topic_id is required")
	}
	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, fmt.Errorf("googlepubsub: failed to create client: %w", err)
	}
	return &GooglePubSubNotifier{projectID: projectID, topicID: topicID, client: client}, nil
}

func newGooglePubSubNotifierWithClient(projectID, topicID string, opts ...option.ClientOption) (*GooglePubSubNotifier, error) {
	if projectID == "" {
		return nil, fmt.Errorf("googlepubsub: project_id is required")
	}
	if topicID == "" {
		return nil, fmt.Errorf("googlepubsub: topic_id is required")
	}
	client, err := pubsub.NewClient(context.Background(), projectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("googlepubsub: failed to create client: %w", err)
	}
	return &GooglePubSubNotifier{projectID: projectID, topicID: topicID, client: client}, nil
}

// Notify publishes a message to the configured Pub/Sub topic.
func (n *GooglePubSubNotifier) Notify(ctx context.Context, secret vault.Secret) error {
	status := "expiring"
	if secret.IsExpired() {
		status = "expired"
	}

	msg := pubSubMessage{
		Path:     secret.Path,
		Status:   status,
		DaysLeft: secret.DaysUntilExpiration(),
		Message:  FormatMessage(secret),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("googlepubsub: failed to marshal message: %w", err)
	}

	topic := n.client.Topic(n.topicID)
	result := topic.Publish(ctx, &pubsub.Message{Data: data})
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("googlepubsub: publish failed: %w", err)
	}
	return nil
}
