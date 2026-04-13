package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// mockDynamoDBClient is a test double for dynamoDBClient.
type mockDynamoDBClient struct {
	calledWith *dynamodb.PutItemInput
	err        error
}

func (m *mockDynamoDBClient) PutItem(_ context.Context, params *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	m.calledWith = params
	return &dynamodb.PutItemOutput{}, m.err
}

func newDynamoDBSecret(daysLeft int) *vault.Secret {
	expiry := time.Now().Add(time.Duration(daysLeft) * 24 * time.Hour)
	return &vault.Secret{Path: "secret/db/password", Expiration: expiry}
}

func TestNewDynamoDBNotifier_MissingTableName(t *testing.T) {
	_, err := newDynamoDBNotifierWithClient(&mockDynamoDBClient{}, "")
	if err == nil {
		t.Fatal("expected error for missing table name")
	}
}

func TestNewDynamoDBNotifier_Valid(t *testing.T) {
	n, err := newDynamoDBNotifierWithClient(&mockDynamoDBClient{}, "vault-events")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.tableName != "vault-events" {
		t.Errorf("expected table name 'vault-events', got %q", n.tableName)
	}
}

func TestDynamoDBNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockDynamoDBClient{}
	n, _ := newDynamoDBNotifierWithClient(mock, "vault-events")

	secret := newDynamoDBSecret(5)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calledWith == nil {
		t.Fatal("expected PutItem to be called")
	}
	if *mock.calledWith.TableName != "vault-events" {
		t.Errorf("unexpected table name: %s", *mock.calledWith.TableName)
	}
	if _, ok := mock.calledWith.Item["SecretPath"]; !ok {
		t.Error("expected SecretPath attribute in item")
	}
}

func TestDynamoDBNotifier_Notify_Expired(t *testing.T) {
	mock := &mockDynamoDBClient{}
	n, _ := newDynamoDBNotifierWithClient(mock, "vault-events")

	secret := newDynamoDBSecret(-1)
	if err := n.Notify(context.Background(), secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calledWith == nil {
		t.Fatal("expected PutItem to be called")
	}
}

func TestDynamoDBNotifier_Notify_ClientError(t *testing.T) {
	mock := &mockDynamoDBClient{err: errors.New("connection refused")}
	n, _ := newDynamoDBNotifierWithClient(mock, "vault-events")

	secret := newDynamoDBSecret(3)
	err := n.Notify(context.Background(), secret)
	if err == nil {
		t.Fatal("expected error from client")
	}
}

func TestDynamoDBNotifier_ImplementsInterface(t *testing.T) {
	mock := &mockDynamoDBClient{}
	n, _ := newDynamoDBNotifierWithClient(mock, "vault-events")
	var _ Notifier = n
}
