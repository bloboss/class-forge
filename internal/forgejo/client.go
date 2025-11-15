package forgejo

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"code.gitea.io/sdk/gitea"
	"go.uber.org/zap"
)

// Client wraps the Gitea SDK client with additional functionality for Forgejo Classroom
type Client struct {
	client  *gitea.Client
	baseURL string
	token   string
	logger  *zap.Logger
	timeout time.Duration
}

// ClientConfig holds configuration for creating a Forgejo client
type ClientConfig struct {
	BaseURL   string
	Token     string
	Timeout   time.Duration
	Logger    *zap.Logger
	UserAgent string
}

// NewClient creates a new Forgejo client with the provided configuration
func NewClient(cfg ClientConfig) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("API token is required")
	}

	// Set defaults
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.Logger == nil {
		cfg.Logger = zap.NewNop()
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "forgejo-classroom/1.0"
	}

	// Create custom HTTP client with timeout
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}

	// Create Gitea SDK client
	client, err := gitea.NewClient(
		cfg.BaseURL,
		gitea.SetToken(cfg.Token),
		gitea.SetHTTPClient(httpClient),
		gitea.SetUserAgent(cfg.UserAgent),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gitea client: %w", err)
	}

	return &Client{
		client:  client,
		baseURL: cfg.BaseURL,
		token:   cfg.Token,
		logger:  cfg.Logger,
		timeout: cfg.Timeout,
	}, nil
}

// HealthCheck verifies connectivity to the Forgejo instance and validates the API token
func (c *Client) HealthCheck(ctx context.Context) error {
	c.logger.Debug("performing Forgejo health check", zap.String("base_url", c.baseURL))

	// Get server version to verify connectivity
	version, _, err := c.client.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to get server version: %w", err)
	}

	c.logger.Info("Forgejo connection successful",
		zap.String("version", version),
		zap.String("base_url", c.baseURL),
	)

	// Verify token by attempting to get current user
	user, _, err := c.client.GetMyUserInfo()
	if err != nil {
		return fmt.Errorf("failed to verify API token: %w", err)
	}

	c.logger.Info("Forgejo API token validated",
		zap.String("username", user.UserName),
		zap.Int64("user_id", user.ID),
	)

	return nil
}

// GetVersion returns the Forgejo server version
func (c *Client) GetVersion(ctx context.Context) (string, error) {
	version, _, err := c.client.ServerVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get server version: %w", err)
	}
	return version, nil
}

// GetCurrentUser returns the authenticated user information
func (c *Client) GetCurrentUser(ctx context.Context) (*gitea.User, error) {
	user, _, err := c.client.GetMyUserInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	return user, nil
}

// Close closes any resources held by the client
func (c *Client) Close() error {
	// The Gitea SDK client doesn't require explicit cleanup
	c.logger.Debug("closing Forgejo client")
	return nil
}
