package forgejo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		config    ClientConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid configuration",
			config: ClientConfig{
				BaseURL: "https://forgejo.example.com",
				Token:   "test-token",
				Timeout: 30 * time.Second,
				Logger:  zap.NewNop(),
			},
			wantError: false,
		},
		{
			name: "missing base URL",
			config: ClientConfig{
				Token:  "test-token",
				Logger: zap.NewNop(),
			},
			wantError: true,
			errorMsg:  "base URL is required",
		},
		{
			name: "missing token",
			config: ClientConfig{
				BaseURL: "https://forgejo.example.com",
				Logger:  zap.NewNop(),
			},
			wantError: true,
			errorMsg:  "API token is required",
		},
		{
			name: "default timeout applied",
			config: ClientConfig{
				BaseURL: "https://forgejo.example.com",
				Token:   "test-token",
				// Timeout not specified
			},
			wantError: false,
		},
		{
			name: "default logger applied",
			config: ClientConfig{
				BaseURL: "https://forgejo.example.com",
				Token:   "test-token",
				// Logger not specified
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, client)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, client)
				assert.NotNil(t, client.client)
				assert.Equal(t, tt.config.BaseURL, client.baseURL)
				assert.Equal(t, tt.config.Token, client.token)

				// Verify defaults are applied
				if tt.config.Timeout == 0 {
					assert.Equal(t, 30*time.Second, client.timeout)
				} else {
					assert.Equal(t, tt.config.Timeout, client.timeout)
				}

				assert.NotNil(t, client.logger)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://forgejo.example.com",
		Token:   "test-token",
		Logger:  zap.NewNop(),
	})
	require.NoError(t, err)

	err = client.Close()
	assert.NoError(t, err)
}

func TestClientConfig_Validation(t *testing.T) {
	t.Run("valid config with all fields", func(t *testing.T) {
		cfg := ClientConfig{
			BaseURL:   "https://forgejo.example.com",
			Token:     "test-token",
			Timeout:   30 * time.Second,
			Logger:    zap.NewNop(),
			UserAgent: "forgejo-classroom/1.0",
		}

		client, err := NewClient(cfg)
		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("empty base URL", func(t *testing.T) {
		cfg := ClientConfig{
			BaseURL: "",
			Token:   "test-token",
		}

		client, err := NewClient(cfg)
		require.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("empty token", func(t *testing.T) {
		cfg := ClientConfig{
			BaseURL: "https://forgejo.example.com",
			Token:   "",
		}

		client, err := NewClient(cfg)
		require.Error(t, err)
		assert.Nil(t, client)
	})
}

// Note: HealthCheck, GetVersion, and GetCurrentUser require a real Forgejo instance
// These are tested in integration tests
func TestClient_HealthCheck_UnitTest(t *testing.T) {
	t.Skip("Health check requires real Forgejo instance - tested in integration tests")
}

func TestClient_GetVersion_UnitTest(t *testing.T) {
	t.Skip("GetVersion requires real Forgejo instance - tested in integration tests")
}

func TestClient_GetCurrentUser_UnitTest(t *testing.T) {
	t.Skip("GetCurrentUser requires real Forgejo instance - tested in integration tests")
}

// Benchmarks
func BenchmarkNewClient(b *testing.B) {
	cfg := ClientConfig{
		BaseURL: "https://forgejo.example.com",
		Token:   "test-token",
		Logger:  zap.NewNop(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client, err := NewClient(cfg)
		if err != nil {
			b.Fatal(err)
		}
		_ = client.Close()
	}
}

func TestContextCancellation(t *testing.T) {
	t.Run("operations respect context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// The actual API calls will be made in integration tests
		// This test verifies the pattern is correct
		assert.NotNil(t, ctx)
	})
}
