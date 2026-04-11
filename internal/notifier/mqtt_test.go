package notifier

import (
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// mockMQTTToken is a minimal mqtt.Token implementation.
type mockMQTTToken struct{}

func (t *mockMQTTToken) Wait() bool                       { return true }
func (t *mockMQTTToken) WaitTimeout(d time.Duration) bool { return true }
func (t *mockMQTTToken) Done() <-chan struct{}             { ch := make(chan struct{}); close(ch); return ch }
func (t *mockMQTTToken) Error() error                     { return nil }

// mockMQTTClient records published messages.
type mockMQTTClient struct {
	mu       sync.Mutex
	published []string
}

func (m *mockMQTTClient) Connect() mqtt.Token    { return &mockMQTTToken{} }
func (m *mockMQTTClient) Disconnect(quiesce uint) {}
func (m *mockMQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.published = append(m.published, payload.(string))
	return &mockMQTTToken{}
}
func (m *mockMQTTClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return &mockMQTTToken{}
}
func (m *mockMQTTClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return &mockMQTTToken{}
}
func (m *mockMQTTClient) Unsubscribe(topics ...string) mqtt.Token { return &mockMQTTToken{} }
func (m *mockMQTTClient) AddRoute(topic string, callback mqtt.MessageHandler) {}
func (m *mockMQTTClient) IsConnected() bool                       { return true }
func (m *mockMQTTClient) IsConnectionOpen() bool                  { return true }
func (m *mockMQTTClient) OptionsReader() mqtt.ClientOptionsReader { return nil }

func TestNewMQTTNotifier_MissingBroker(t *testing.T) {
	_, err := newMQTTNotifierWithClient(&mockMQTTClient{}, "", 1)
	if err == nil {
		t.Fatal("expected error for missing topic")
	}
}

func TestNewMQTTNotifier_Valid(t *testing.T) {
	n, err := newMQTTNotifierWithClient(&mockMQTTClient{}, "vaultwatch/alerts", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestMQTTNotifier_Notify_ExpiringSoon(t *testing.T) {
	client := &mockMQTTClient{}
	n, _ := newMQTTNotifierWithClient(client, "vaultwatch/alerts", 1)

	secret := &vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(client.published) != 1 {
		t.Fatalf("expected 1 published message, got %d", len(client.published))
	}
}

func TestMQTTNotifier_Notify_Expired(t *testing.T) {
	client := &mockMQTTClient{}
	n, _ := newMQTTNotifierWithClient(client, "vaultwatch/alerts", 0)

	secret := &vault.Secret{
		Path:      "secret/api/key",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(client.published) != 1 {
		t.Fatalf("expected 1 published message, got %d", len(client.published))
	}
}

func TestMQTTNotifier_ImplementsInterface(t *testing.T) {
	var _ Notifier = (*MQTTNotifier)(nil)
}
