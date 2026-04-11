package notifier

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// MQTTNotifier publishes secret expiration alerts to an MQTT broker topic.
type MQTTNotifier struct {
	client mqtt.Client
	topic  string
	qos    byte
}

// NewMQTTNotifier creates a new MQTTNotifier.
// brokerURL should be in the form "tcp://host:port".
func NewMQTTNotifier(brokerURL, topic, clientID string, qos byte) (*MQTTNotifier, error) {
	if brokerURL == "" {
		return nil, fmt.Errorf("mqtt: broker URL is required")
	}
	if topic == "" {
		return nil, fmt.Errorf("mqtt: topic is required")
	}
	if clientID == "" {
		clientID = fmt.Sprintf("vaultwatch-%d", time.Now().UnixNano())
	}
	if qos > 2 {
		qos = 1
	}

	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetConnectTimeout(10 * time.Second).
		SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt: failed to connect to broker: %w", token.Error())
	}

	return &MQTTNotifier{
		client: client,
		topic:  topic,
		qos:    qos,
	}, nil
}

// newMQTTNotifierWithClient creates an MQTTNotifier with an injected client (for testing).
func newMQTTNotifierWithClient(client mqtt.Client, topic string, qos byte) (*MQTTNotifier, error) {
	if topic == "" {
		return nil, fmt.Errorf("mqtt: topic is required")
	}
	return &MQTTNotifier{client: client, topic: topic, qos: qos}, nil
}

// Notify publishes a formatted alert message to the configured MQTT topic.
func (n *MQTTNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	payload := fmt.Sprintf(`{"subject":%q,"body":%q}`, msg.Subject, msg.Body)

	token := n.client.Publish(n.topic, n.qos, false, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt: publish failed: %w", err)
	}
	return nil
}
