package vault

import (
	"context"
	"fmt"
	"time"
)

// Secret represents a Vault secret with expiration metadata
type Secret struct {
	Path           string
	Version        int
	CreatedTime    time.Time
	ExpirationTime *time.Time
	Metadata       map[string]interface{}
}

// DaysUntilExpiration returns the number of days until the secret expires
// Returns -1 if the secret has no expiration
func (s *Secret) DaysUntilExpiration() int {
	if s.ExpirationTime == nil {
		return -1
	}

	duration := time.Until(*s.ExpirationTime)
	days := int(duration.Hours() / 24)
	return days
}

// IsExpired checks if the secret has expired
func (s *Secret) IsExpired() bool {
	if s.ExpirationTime == nil {
		return false
	}
	return time.Now().After(*s.ExpirationTime)
}

// IsExpiringSoon checks if the secret will expire within the given threshold
func (s *Secret) IsExpiringSoon(threshold time.Duration) bool {
	if s.ExpirationTime == nil {
		return false
	}
	return time.Until(*s.ExpirationTime) <= threshold
}

// GetSecret retrieves a secret from Vault at the specified path
func (c *Client) GetSecret(ctx context.Context, path string) (*Secret, error) {
	if path == "" {
		return nil, fmt.Errorf("secret path cannot be empty")
	}

	secretData, err := c.api.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %s: %w", path, err)
	}

	if secretData == nil {
		return nil, fmt.Errorf("secret not found at path: %s", path)
	}

	secret := &Secret{
		Path:     path,
		Metadata: secretData.Data,
	}

	// Extract metadata if available (KV v2)
	if metadata, ok := secretData.Data["metadata"].(map[string]interface{}); ok {
		if createdTime, ok := metadata["created_time"].(string); ok {
			if t, err := time.Parse(time.RFC3339, createdTime); err == nil {
				secret.CreatedTime = t
			}
		}
	}

	// Check for lease duration (dynamic secrets)
	if secretData.LeaseDuration > 0 {
		expTime := time.Now().Add(time.Duration(secretData.LeaseDuration) * time.Second)
		secret.ExpirationTime = &expTime
	}

	return secret, nil
}
