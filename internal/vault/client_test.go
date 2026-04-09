package vault

import (
	"context"
	"testing"
	"time"
)

func TestNewClient_ValidConfig(t *testing.T) {
	cfg := &Config{
		Address: "http://localhost:8200",
		Token:   "test-token",
		Timeout: 10 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("expected client to be non-nil")
	}

	if client.API() == nil {
		t.Fatal("expected API client to be non-nil")
	}
}

func TestNewClient_MissingAddress(t *testing.T) {
	cfg := &Config{
		Token: "test-token",
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	cfg := &Config{
		Address: "http://localhost:8200",
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewClient_WithNamespace(t *testing.T) {
	cfg := &Config{
		Address:   "http://localhost:8200",
		Token:     "test-token",
		Namespace: "test-namespace",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("expected client to be non-nil")
	}
}

func TestNewClient_DefaultTimeout(t *testing.T) {
	cfg := &Config{
		Address: "http://localhost:8200",
		Token:   "test-token",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client.API().Timeout() != 30*time.Second {
		t.Errorf("expected default timeout of 30s, got %v", client.API().Timeout())
	}
}

func TestHealthCheck_Context(t *testing.T) {
	cfg := &Config{
		Address: "http://localhost:8200",
		Token:   "test-token",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// This will fail in tests without a real Vault, but validates the method signature
	err = client.HealthCheck(ctx)
	// We expect an error since there's no real Vault running
	if err == nil {
		t.Log("health check succeeded (vault must be running)")
	}
}
