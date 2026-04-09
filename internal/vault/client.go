package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with custom methods
type Client struct {
	api *vaultapi.Client
}

// Config holds Vault client configuration
type Config struct {
	Address   string
	Token     string
	Namespace string
	Timeout   time.Duration
}

// NewClient creates a new Vault client with the provided configuration
func NewClient(cfg *Config) (*Client, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("vault address is required")
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("vault token is required")
	}

	config := vaultapi.DefaultConfig()
	config.Address = cfg.Address

	if cfg.Timeout > 0 {
		config.Timeout = cfg.Timeout
	} else {
		config.Timeout = 30 * time.Second
	}

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	client.SetToken(cfg.Token)

	if cfg.Namespace != "" {
		client.SetNamespace(cfg.Namespace)
	}

	return &Client{api: client}, nil
}

// HealthCheck verifies connectivity to Vault
func (c *Client) HealthCheck(ctx context.Context) error {
	resp, err := c.api.Sys().HealthWithContext(ctx)
	if err != nil {
		return fmt.Errorf("vault health check failed: %w", err)
	}

	if !resp.Initialized {
		return fmt.Errorf("vault is not initialized")
	}

	if resp.Sealed {
		return fmt.Errorf("vault is sealed")
	}

	return nil
}

// API returns the underlying Vault API client
func (c *Client) API() *vaultapi.Client {
	return c.api
}
