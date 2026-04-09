package vault

import (
	"testing"
	"time"
)

func TestSecret_DaysUntilExpiration(t *testing.T) {
	tests := []struct {
		name     string
		expTime  *time.Time
		expected int
	}{
		{
			name:     "no expiration",
			expTime:  nil,
			expected: -1,
		},
		{
			name: "expires in 5 days",
			expTime: func() *time.Time {
				t := time.Now().Add(5 * 24 * time.Hour)
				return &t
			}(),
			expected: 5,
		},
		{
			name: "expires in 30 days",
			expTime: func() *time.Time {
				t := time.Now().Add(30 * 24 * time.Hour)
				return &t
			}(),
			expected: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Secret{ExpirationTime: tt.expTime}
			result := s.DaysUntilExpiration()
			if result != tt.expected {
				t.Errorf("expected %d days, got %d", tt.expected, result)
			}
		})
	}
}

func TestSecret_IsExpired(t *testing.T) {
	pastTime := time.Now().Add(-1 * time.Hour)
	futureTime := time.Now().Add(1 * time.Hour)

	tests := []struct {
		name     string
		expTime  *time.Time
		expected bool
	}{
		{"no expiration", nil, false},
		{"expired", &pastTime, true},
		{"not expired", &futureTime, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Secret{ExpirationTime: tt.expTime}
			if s.IsExpired() != tt.expected {
				t.Errorf("expected IsExpired() = %v, got %v", tt.expected, s.IsExpired())
			}
		})
	}
}

func TestSecret_IsExpiringSoon(t *testing.T) {
	threshold := 7 * 24 * time.Hour // 7 days

	tests := []struct {
		name     string
		expTime  *time.Time
		expected bool
	}{
		{
			name:     "no expiration",
			expTime:  nil,
			expected: false,
		},
		{
			name: "expires in 3 days",
			expTime: func() *time.Time {
				t := time.Now().Add(3 * 24 * time.Hour)
				return &t
			}(),
			expected: true,
		},
		{
			name: "expires in 30 days",
			expTime: func() *time.Time {
				t := time.Now().Add(30 * 24 * time.Hour)
				return &t
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Secret{ExpirationTime: tt.expTime}
			if s.IsExpiringSoon(threshold) != tt.expected {
				t.Errorf("expected IsExpiringSoon() = %v, got %v", tt.expected, s.IsExpiringSoon(threshold))
			}
		})
	}
}
